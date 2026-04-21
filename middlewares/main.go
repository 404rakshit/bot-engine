package middleware

import authMiddlewares "bot-engine/middlewares/auth"

type Middlerware struct {
	AuthMiddleware authMiddlewares.AuthMiddleware
}

func NewMiddlerware() *Middlerware {
	return &Middlerware{
		AuthMiddleware: authMiddlewares.NewAuthMiddleware(),
	}
}
