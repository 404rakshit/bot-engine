package middleware

import authMiddlewares "di/middlewares/auth"

type Middlerware struct {
	AuthMiddleware authMiddlewares.AuthMiddleware
}

func NewMiddlerware() *Middlerware {
	return &Middlerware{
		AuthMiddleware: authMiddlewares.NewAuthMiddleware(),
	}
}
