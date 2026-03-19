package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Bot struct {
	token   string
	baseURL string
	client  *http.Client
	handler CommandHandler
	offset  int
}

type CommandHandler interface {
	Handle(ctx context.Context, update Update) error
}

func New(token string, handler CommandHandler) *Bot {
	return &Bot{
		token:   token,
		baseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		client:  &http.Client{Timeout: 30 * time.Second},
		handler: handler,
	}
}

// StartPolling runs long-polling loop to receive updates.
func (b *Bot) StartPolling(ctx context.Context) error {
	log.Println("Bot started with long-polling")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			updates, err := b.getUpdates(ctx)
			if err != nil {
				log.Printf("getUpdates error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			for _, update := range updates {
				if err := b.handler.Handle(ctx, update); err != nil {
					log.Printf("handler error: %v", err)
				}
				b.offset = update.UpdateID + 1
			}
		}
	}
}

// HandleWebhook returns an http.HandlerFunc for webhook mode.
func (b *Bot) HandleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if err := b.handler.Handle(r.Context(), update); err != nil {
			log.Printf("handler error: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (b *Bot) getUpdates(ctx context.Context) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", b.baseURL, b.offset)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("getUpdates: build request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getUpdates: http call: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("getUpdates: decode response: %w", err)
	}

	return result.Result, nil
}

// SendMessage sends a text message to a chat.
func (b *Bot) SendMessage(ctx context.Context, chatID int64, text string) error {
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/sendMessage", b.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sendMessage: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("sendMessage: http call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sendMessage: telegram API error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
