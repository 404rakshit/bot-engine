package middleware

import (
	"bot-engine/config"
	"bot-engine/middlewares/auth"
	authMiddlewares "bot-engine/middlewares/auth"
)

type Middlewares struct {
	AuthMiddleware authMiddlewares.AuthMiddleware
}

func NewMiddlewares(cfg *config.AuthSecretsConfig) *Middlewares {
	return &Middlewares{
		// Pass the JWT secret from your config
		AuthMiddleware: auth.NewAuthMiddleware(cfg.JWTSecret),
	}
}
