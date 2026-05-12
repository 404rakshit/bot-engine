// services/bot/webhook_registrar.go
package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookRegistrar interface {
	Register(token string, webhookURL string) error
}

type webhookRegistrar struct {
	client *http.Client
}

func NewWebhookRegistrar() WebhookRegistrar {
	return &webhookRegistrar{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *webhookRegistrar) Register(token string, webhookURL string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", token)

	payload := map[string]string{
		"url": webhookURL,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := r.client.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to reach Telegram API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram returned status %d", resp.StatusCode)
	}

	return nil
}
