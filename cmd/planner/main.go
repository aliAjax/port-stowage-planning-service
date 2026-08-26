package main

import (
	"context"
	"errors"
	"github.com/example/port-stowage-planner/internal/config"
	"github.com/example/port-stowage-planner/internal/dispatch"
	"github.com/example/port-stowage-planner/internal/repository"
	"github.com/example/port-stowage-planner/internal/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	srv := transport.New(repository.New(), dispatch.NewSimulator())
	httpSrv := &http.Server{Addr: cfg.Addr, Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go httpSrv.ListenAndServe()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdown); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
