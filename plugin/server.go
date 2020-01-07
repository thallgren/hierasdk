//+build !test

// Package plugin exposes the API for starting the RESTful plugin service.
package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/routes"
)

const defaultMinPort = 10000
const defaultMaxPort = 25000

// ServeAndExit starts serving the plug-in
func ServeAndExit() {
	minPort := getEnvInt(`HIERA_MIN_PORT`, defaultMinPort)
	maxPort := getEnvInt(`HIERA_MAX_PORT`, defaultMaxPort)
	os.Exit(Serve(os.Args[0], minPort, maxPort, os.Stdout, os.Stderr))
}

// Serve starts serving the plug-in using the given name, port range, stderr, and stdout
func Serve(name string, minPort, maxPort int, stdout, stderr io.Writer) int {
	if getEnvInt(`HIERA_MAGIC_COOKIE`, 0) != hiera.MagicCookie {
		_, _ = fmt.Fprintf(stderr,
			"%s is meant to be used as a Hiera RESTful plugin. It should not be started from a command shell\n", name)
		return 1
	}
	if minPort > maxPort {
		_, _ = fmt.Fprintf(os.Stderr, "min port %d is greater than max port %d\n", minPort, maxPort)
		return 1
	}

	sockDir := os.Getenv("HIERA_PLUGIN_SOCKET_DIR")

	var listener net.Listener
	var err error
	if sockDir == "" {
		listener, err = getTCPListener(minPort, maxPort)
	} else {
		listener, err = getSocketListener(sockDir, path.Base(name))
	}
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	handler, functions := routes.Register()
	return startServer(listener, handler, functions, stdout, stderr)
}

func getTCPListener(minPort, maxPort int) (net.Listener, error) {
	for port := minPort; port <= maxPort; port++ {
		listener, err := net.Listen(`tcp`, `127.0.0.1:`+strconv.Itoa(port))
		if err == nil {
			return listener, nil
		}
	}
	return nil, fmt.Errorf(`no available port in the range %d to %d`, minPort, maxPort)
}

var tempFileAttempts = 10

// tempFile generates random file name in a given directory
// file name has a form of <plugin-name><random-string>.socket
// the function returns an error is given directory doesn't exist.
func tempFile(dir, prefix string) (string, error) {
	fi, err := os.Lstat(dir)
	if err != nil || !fi.IsDir() {
		return "", fmt.Errorf("path is not a directory %s", dir)
	}

	filename := prefix
	seed := rand.New(rand.NewSource(int64(time.Now().UnixNano() + int64(os.Getpid()))))

	for i := 0; i < tempFileAttempts; i++ {
		path := filepath.Join(dir, filename+".socket")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path, nil
		}

		filename = filename + strconv.Itoa(seed.Int())[:1]
	}

	return "", fmt.Errorf("failed to generate temporary file in %s", dir)
}

func getSocketListener(dir, name string) (net.Listener, error) {
	socket, err := tempFile(dir, name)
	if err != nil {
		return nil, err
	}

	return net.Listen(`unix`, socket)
}

func getEnvInt(n string, defaultValue int) int {
	if v := os.Getenv(n); len(v) > 0 {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

func startServer(listener net.Listener, router http.Handler, functions dgo.Map, ow, ew io.Writer) int {
	err := json.NewEncoder(ow).Encode(vf.Map(`version`, hiera.ProtoVersion, `network`, listener.Addr().Network(), `address`, listener.Addr().String(), `functions`, functions))
	if err != nil {
		_, _ = fmt.Fprintln(ew, err)
		return 1
	}

	server := http.Server{Handler: router}
	done := make(chan bool, 1)
	// Allow graceful shutdown of server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			_, _ = fmt.Fprintf(ew, "Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	if err = server.Serve(listener); err != nil && err != http.ErrServerClosed {
		_, _ = fmt.Fprintf(ew, "Could not listen on %s: %v\n", listener.Addr(), err)
		return 1
	}
	<-done
	return 0
}
