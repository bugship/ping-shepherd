package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bugship/ping-shepherd/internal/config"
	"github.com/bugship/ping-shepherd/internal/httpapi"
	"github.com/bugship/ping-shepherd/internal/store"
)

func main() {
	cfg := config.FromEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := store.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("postgres not ready (%v); /ready will fail until it is", err)
	} else if err := db.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if db != nil {
		defer db.Close()
	}

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           httpapi.New(db),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("ping-shepherd listening on %s", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
