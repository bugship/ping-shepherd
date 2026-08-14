// Package notify sends alert messages.
package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Telegram sends text to a chat via the Bot API.
type Telegram struct {
	Token  string
	ChatID string
	Client *http.Client
	Base   string
}

// Send posts text to the configured chat.
func (t Telegram) Send(ctx context.Context, text string) error {
	if strings.TrimSpace(t.Token) == "" || strings.TrimSpace(t.ChatID) == "" {
		return fmt.Errorf("telegram is not configured")
	}
	base := t.Base
	if base == "" {
		base = "https://api.telegram.org"
	}
	client := t.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	body, err := json.Marshal(map[string]any{
		"chat_id":                  t.ChatID,
		"text":                     text,
		"disable_web_page_preview": true,
	})
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/bot%s/sendMessage", strings.TrimRight(base, "/"), t.Token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return fmt.Errorf("telegram status %d", res.StatusCode)
	}
	return nil
}
