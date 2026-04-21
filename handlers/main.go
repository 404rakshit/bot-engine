package handlers

import (
	"di/services"

	userHandlers "di/handlers/users"
)

type Handlers struct {
	UserHandler *userHandlers.UserHandler
}

func NewHandler(
	s *services.Services,
) *Handlers {
	return &Handlers{
		UserHandler: userHandlers.NewUserHandler(s.UserService),
	}
}
