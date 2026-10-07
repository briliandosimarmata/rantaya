package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"ruang.local/api/internal/app"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	c := app.Configuration()
	a, err := app.New(ctx, c)
	if err != nil {
		slog.Error("API startup failed", "error", err)
		os.Exit(1)
	}
	defer a.DB.Close()
	server := &http.Server{Addr: c.Addr, Handler: a.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 25 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 90 * time.Second, MaxHeaderBytes: 1 << 20}
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := a.Remind(ctx); err != nil {
					slog.Error("event reminders", "error", err)
				}
				if err := a.Expire(ctx); err != nil {
					slog.Error("expire orders", "error", err)
				}
			}
		}
	}()
	go func() {
		<-ctx.Done()
		shutdown, done := context.WithTimeout(context.Background(), 10*time.Second)
		defer done()
		_ = server.Shutdown(shutdown)
	}()
	slog.Info("Go API ready", "address", c.Addr, "demo", c.Demo)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
