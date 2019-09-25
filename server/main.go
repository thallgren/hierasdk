//+build !test

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/lyraproj/hierasdk/routes"
)

func main() {
	os.Exit(run(os.Args, os.Stderr))
}

func run(args []string, errStr io.Writer) int {
	port := ``
	if len(args) == 2 {
		port = args[1]
		if _, err := strconv.Atoi(port); err != nil {
			port = ``
		}
	}
	if port == `` {
		_, _ = fmt.Fprintf(errStr, "usage: %s <port number>", args[0])
		return 1
	}
	return startServer(routes.Register(), `:`+port, errStr)
}

func startServer(router *http.ServeMux, listenAddr string, errStr io.Writer) int {
	server := http.Server{Addr: listenAddr, Handler: router}
	done := make(chan bool, 1)
	// Allow graceful shutdown of server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			_, _ = fmt.Fprintf(errStr, "Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		_, _ = fmt.Fprintf(errStr, "Could not listen on %s: %v\n", listenAddr, err)
		return 1
	}
	<-done
	return 0
}
