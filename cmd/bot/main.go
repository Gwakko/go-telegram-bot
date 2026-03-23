package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gwakko/go-telegram-bot/internal/bot"
	"github.com/gwakko/go-telegram-bot/internal/handlers"
	"github.com/gwakko/go-telegram-bot/internal/storage"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	mode := os.Getenv("BOT_MODE") // "polling" or "webhook"
	if mode == "" {
		mode = "polling"
	}

	store := storage.NewMemoryStore()
	defer store.Close()

	// Bot is created first, router wraps it with command handling
	var b *bot.Bot
	router := handlers.NewRouter(nil, store) // temporary nil, set below

	b = bot.New(token, router)
	// Re-create router with actual bot reference
	router = handlers.NewRouter(b, store)
	b = bot.New(token, router)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	switch mode {
	case "webhook":
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/webhook", b.HandleWebhook())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		})

		log.Printf("Webhook server starting on :%s", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
			log.Fatalf("server error: %v", err)
		}

	default:
		log.Println("Starting in polling mode")
		if err := b.StartPolling(ctx); err != nil && ctx.Err() == nil {
			log.Fatalf("polling error: %v", err)
		}
	}
}
