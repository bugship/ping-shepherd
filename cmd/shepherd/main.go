// Command shepherd runs the Ping Shepherd server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bugship/ping-shepherd/internal/checker"
	"github.com/bugship/ping-shepherd/internal/config"
	"github.com/bugship/ping-shepherd/internal/httpapi"
	"github.com/bugship/ping-shepherd/internal/notify"
	"github.com/bugship/ping-shepherd/internal/store"
)

func main() {
	cfg := config.FromEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	backend, closer := openStore(ctx, cfg.DatabaseURL)
	if closer != nil {
		defer closer()
	}

	var alert checker.Sender
	if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
		alert = notify.Telegram{Token: cfg.TelegramToken, ChatID: cfg.TelegramChatID}
		log.Printf("telegram alerts enabled")
	}
	go checker.Loop(ctx, backend, cfg.CheckInterval, cfg.CheckTimeout, alert)

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           httpapi.New(backend),
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

func openStore(ctx context.Context, databaseURL string) (httpapi.Backend, func()) {
	db, err := store.Open(ctx, databaseURL)
	if err != nil {
		log.Printf("postgres unavailable (%v); using memory store", err)
		return store.NewMemory(), nil
	}
	if err := db.Migrate(ctx); err != nil {
		log.Printf("migrate failed (%v); using memory store", err)
		db.Close()
		return store.NewMemory(), nil
	}
	return db, db.Close
}
