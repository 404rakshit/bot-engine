package handlers

import (
	"bot-engine/config"
	"bot-engine/services"

	authHandlers "bot-engine/handlers/auth"
	userHandlers "bot-engine/handlers/users"
)

type Handlers struct {
	UserHandler *userHandlers.UserHandler
	AuthHandler *authHandlers.AuthHandler
}

func NewHandler(
	s *services.Services,
	botCfg *config.BotSecretsConfig,
) *Handlers {
	return &Handlers{
		UserHandler: userHandlers.NewUserHandler(s.UserService, s.BotService),
		AuthHandler: authHandlers.NewAuthHandler(s.AuthService, botCfg),
	}
}
