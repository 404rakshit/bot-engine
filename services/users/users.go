package users

import (
	"context"
	userModels "di/models/mongo/users"
	userRepos "di/repositories/users"
)

type userRepo = userRepos.UserRepository
type userModel = userModels.User

type UserService interface {
	List(ctx context.Context) ([]userModel, error)
	Create(ctx context.Context, data *userModel) error
}

type userService struct {
	userRepo userRepo
}

func NewUserService(userRepo userRepo) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) List(ctx context.Context) ([]userModel, error) {
	users, err := s.userRepo.List(ctx)
	return users, err
}

func (s *userService) Create(ctx context.Context, data *userModel) error {

	if err := data.Validate(); err != nil {
		return err
	}

	if err := s.userRepo.Create(ctx, data); err != nil {
		return err
	}

	return nil
}
