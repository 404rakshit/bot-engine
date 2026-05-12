package handlers

import (
	"bot-engine/config"
	"bot-engine/services"

	authHandlers "bot-engine/handlers/auth"
	webhookHandlers "bot-engine/handlers/engine"
	userHandlers "bot-engine/handlers/users"
)

type Handlers struct {
	UserHandler    *userHandlers.UserHandler
	AuthHandler    *authHandlers.AuthHandler
	WebhookHandler *webhookHandlers.WebhookHandler
}

func NewHandler(
	s *services.Services,
	botCfg *config.BotSecretsConfig,
) *Handlers {
	return &Handlers{
		UserHandler:    userHandlers.NewUserHandler(s.UserService, s.BotService),
		AuthHandler:    authHandlers.NewAuthHandler(s.AuthService, botCfg),
		WebhookHandler: webhookHandlers.NewWebhookHandler(s.BotService, s.EngineService),
	}
}
