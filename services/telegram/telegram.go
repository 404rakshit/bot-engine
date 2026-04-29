package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TelegramService interface {
	VerifyToken(ctx context.Context, token string) (username string, botID int64, err error)
}

type telegramService struct {
	client *http.Client
}

func NewTelegramService() *telegramService {
	return &telegramService{
		client: &http.Client{},
	}
}

type telegramResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	Result      struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"result"`
}

func (s *telegramService) VerifyToken(ctx context.Context, token string) (username string, botID int64, err error) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)

	if err != nil {
		return "", 0, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	var tgResp telegramResponse

	if err := json.NewDecoder(resp.Body).Decode(&tgResp); err != nil {
		return "", 0, fmt.Errorf("failed to decode telegram response: %w", err)
	}

	if !tgResp.Ok {
		// If Telegram returns a 401 Unauthorized, the token is invalid
		if resp.StatusCode == http.StatusUnauthorized {
			return "", 0, errors.New("invalid telegram bot token")
		}
		// Catch-all for other Telegram API errors
		return "", 0, fmt.Errorf("telegram API error: %s", tgResp.Description)
	}

	return tgResp.Result.Username, tgResp.Result.ID, nil
}
