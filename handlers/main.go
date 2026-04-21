package handlers

import (
	"bot-engine/services"

	userHandlers "bot-engine/handlers/users"
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
