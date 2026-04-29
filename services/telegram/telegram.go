package telegram

import "context"

type TelegramService interface {
	VerifyToken(ctx context.Context, token string) (username string, botID int64, err error)
}

type telegramService struct {
}

func NewTelegramService() *telegramService {
	return &telegramService{}
}

func (*telegramService) VerifyToken(ctx context.Context, token string) (username string, botID int64, err error) {
	return "", 0, nil
}
