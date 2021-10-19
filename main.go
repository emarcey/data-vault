package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"emarcey/data-vault/dependencies"
	"emarcey/data-vault/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts, err := dependencies.ReadOpts("./server_conf.yml")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	deps, err := dependencies.MakeDependencies(ctx, opts)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	service := server.NewService(opts.Version, deps)
	handler := server.MakeHttpHandler(service, deps)

	// Listen for application termination.
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	readTimeout, _ := time.ParseDuration("30s")
	writTimeout, _ := time.ParseDuration("30s")
	idleTimeout, _ := time.ParseDuration("60s")

	fmt.Println(opts)
	serve := &http.Server{
		Addr:         opts.HttpAddr,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writTimeout,
		IdleTimeout:  idleTimeout,
	}
	shutdownServer := func() {
		if err := serve.Shutdown(context.TODO()); err != nil {
			deps.Logger.Info("shutdown", err)
		}
	}

	// Start main HTTP server
	go func() {
		deps.Logger.Infof(opts.HttpAddr)

		deps.Logger.Infof(fmt.Sprintf("startup binding to %s for HTTP server", opts.HttpAddr))
		if err := serve.ListenAndServe(); err != nil {
			errs <- err
			deps.Logger.Info("exit ", err)
		}
	}()

	if err := <-errs; err != nil {
		shutdownServer()
		deps.Logger.Info("exit ", err)
	}
}
