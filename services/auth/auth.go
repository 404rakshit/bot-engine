package auth

import (
	"context"
	"errors"
	"time"

	userModels "bot-engine/models/mongo/users"
	userRepos "bot-engine/repositories/users"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Define the contract
type AuthService interface {
	Signup(ctx context.Context, name, email, password string) (*userModels.User, error)
	Login(ctx context.Context, email, password string) (string, *userModels.User, error)
}

type authService struct {
	userRepo  userRepos.UserRepository
	jwtSecret []byte
}

// NewAuthService injects the existing user repository and a JWT secret key (from your config)
func NewAuthService(userRepo userRepos.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Signup(ctx context.Context, name, email, password string) (*userModels.User, error) {
	// 1. Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("a user with this email already exists")
	}

	// 2. Hash the password using bcrypt (Cost of 10 is a standard, safe balance of speed/security)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, errors.New("failed to secure password")
	}

	// 3. Construct the Domain Model
	newUser := &userModels.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	// 4. Validate the model
	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	// 5. Save using your existing User Repo
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create user account")
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *userModels.User, error) {
	// 1. Fetch the user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid email or password") // Keep error generic for security
	}

	// 2. Compare the stored hash with the provided plain-text password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// 3. Generate the JWT Token
	// We store the user's string ID in the "sub" (subject) claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":   time.Now().Unix(),                     // Issued at
	})

	// 4. Sign the token using your secret key
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, errors.New("failed to generate authentication token")
	}

	return tokenString, user, nil
}
