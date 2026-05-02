package auth

import (
	"context"
	"net/http"

	userModels "bot-engine/models/mongo/users"
	"bot-engine/utils"

	"github.com/gin-gonic/gin"
)

// 1. Define DTOs (Data Transfer Objects) for strict validation
type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 2. Define the contract for the Service Layer
// This keeps the handler completely decoupled from bcrypt and JWT logic
type AuthService interface {
	Signup(ctx context.Context, name, email, password string) (*userModels.User, error)
	Login(ctx context.Context, email, password string) (string, *userModels.User, error) // Returns a JWT token and the User
}

type AuthHandler struct {
	AuthService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// SignupUser handles creating a new account
func (h *AuthHandler) SignupUser(c *gin.Context) {
	var req SignupRequest

	// Validate incoming JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response := utils.ErrorResponse("Invalid payload: please check your email and password")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Pass to service layer to hash password and save to DB
	user, err := h.AuthService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response := utils.ErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.SuccessResponse(user, "User signed up successfully")
	c.JSON(http.StatusCreated, response)
}

// LoginUser handles authentication and returning a JWT token
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req LoginRequest

	// Validate incoming JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response := utils.ErrorResponse("Invalid payload: email and password are required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Pass to service layer to verify password and generate JWT
	token, user, err := h.AuthService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// Use 401 Unauthorized for bad credentials
		response := utils.ErrorResponse("Invalid email or password")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	// Bundle the token and user data into a single response map
	data := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	response := utils.SuccessResponse(data, "Login successful")
	c.JSON(http.StatusOK, response)
}
