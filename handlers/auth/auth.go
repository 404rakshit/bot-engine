package auth

import (
	"net/http"
	"os"

	"bot-engine/config"
	authServices "bot-engine/services/auth"
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

type AuthService = authServices.AuthService

type AuthHandler struct {
	AuthService AuthService
	BotCfg      *config.BotSecretsConfig
}

func NewAuthHandler(authService AuthService, botCfg *config.BotSecretsConfig) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		BotCfg:      botCfg,
	}
}

type CustomeUserResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func setAuthCookie(c *gin.Context, token string) {
	// isProduction := os.Getenv("GO_ENV") == "production"

	// Construct the cookie explicitly using Go's standard library
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		MaxAge:   259200,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   true,                  // MUST be true
		SameSite: http.SameSiteNoneMode, // MUST be None
	}

	http.SetCookie(c.Writer, cookie)

	// // Explicitly set the SameSite mode
	// if isProduction {
	// 	cookie.SameSite = http.SameSiteNoneMode
	// } else {
	// 	cookie.SameSite = http.SameSiteLaxMode
	// }

	// // Write the cookie to the response
	// http.SetCookie(c.Writer, cookie)
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
	token, user, err := h.AuthService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response := utils.ErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	setAuthCookie(c, token)

	customResponse := CustomeUserResponse{
		UserID: user.ID.Hex(),
		Email:  user.Email,
	}

	response := utils.SuccessResponse(customResponse, "User signed up successfully")
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

	setAuthCookie(c, token)

	customResponse := CustomeUserResponse{
		UserID: user.ID.Hex(),
		Email:  user.Email,
	}

	response := utils.SuccessResponse(customResponse, "Login successful")
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	isProduction := os.Getenv("GO_ENV") == "production"

	if isProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	// MaxAge -1 deletes the cookie
	c.SetCookie("auth_token", "", -1, "/", "", isProduction, true)

	response := utils.SuccessResponse(nil, "Logged out successfully")
	c.JSON(http.StatusOK, response)
}
