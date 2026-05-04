package auth

import (
	"context"
	"errors"

	"bot-engine/helper"
	userModels "bot-engine/models/mongo/users"

	userRepos "bot-engine/repositories/users"

	"golang.org/x/crypto/bcrypt"
)

// 1. Updated the contract to return the token string for Signup
type AuthService interface {
	Signup(ctx context.Context, name, email, password string) (string, *userModels.User, error)
	Login(ctx context.Context, email, password string) (string, *userModels.User, error)
	ProcessTelegramLogin(ctx context.Context, tgPayload TelegramAuthPayload) (string, error, string)
}

type authService struct {
	userRepo     userRepos.UserRepository
	identityRepo userRepos.IdentityRepository
	jwtSecret    []byte
}

func NewAuthService(
	userRepo userRepos.UserRepository,
	identityRepo userRepos.IdentityRepository,
	jwtSecret string,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		identityRepo: identityRepo,
		jwtSecret:    []byte(jwtSecret),
	}
}

func (s *authService) Signup(
	ctx context.Context,
	name, email,
	password string,
) (
	string,
	*userModels.User,
	error,
) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return "", nil, errors.New("a user with this email already exists")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", nil, errors.New("failed to secure password")
	}

	newUser := &userModels.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	if err := newUser.Validate(); err != nil {
		return "", nil, err
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return "", nil, errors.New("failed to create user account")
	}

	// 2. Generate the token right after the user is saved (so they have an ID)
	tokenString, err := helper.GenerateToken(newUser.ID.Hex(), newUser.Email, s.jwtSecret)

	if err != nil {
		return "", nil, errors.New("user created, but failed to generate token")
	}

	return tokenString, newUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *userModels.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	tokenString, err := helper.GenerateToken(user.ID.Hex(), user.Email, s.jwtSecret)

	if err != nil {
		return "", nil, errors.New("failed to generate authentication token")
	}
	return tokenString, user, nil
}
