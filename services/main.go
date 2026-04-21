package services

import (
	"di/repositories"
	"di/services/users"
)

type Services struct {
	UserService users.UserService
}

func NewServices(
	r *repositories.Repositories,
) *Services {
	return &Services{
		UserService: users.NewUserService(r.UserRepository),
	}
}
