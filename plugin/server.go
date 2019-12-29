//+build !test

// Package plugin exposes the API for starting the RESTful plugin service.
package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
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
	listener, err := getTCPListener(minPort, maxPort)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	handler, functions := routes.Register()
	return startServer(listener, handler, functions, stdout, stderr)
}

func getTCPListener(minPort, maxPort int) (net.Listener, error) {
	var listener net.Listener
	var err error

	sock := os.Getenv("HIERA_PLUGIN_SOCKET")

	for port := minPort; port <= maxPort; port++ {
		if sock == "" {
			listener, err = net.Listen(`tcp`, `127.0.0.1:`+strconv.Itoa(port))
		} else {
			listener, err = net.Listen(`unix`, sock)
		}
		if err == nil {
			return listener, nil
		}
	}
	return nil, fmt.Errorf(`no available port in the range %d to %d`, minPort, maxPort)
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
