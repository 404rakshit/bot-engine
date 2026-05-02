package auth

import (
	"net/http"

	"bot-engine/helper"
	"bot-engine/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	ValidateRequest() gin.HandlerFunc
}

type authMiddleware struct {
	jwtSecret []byte
}

// NewAuthMiddleware requires the JWT secret to verify the incoming tokens.
// You will inject this via Wire using your AuthSecretsConfig.
func NewAuthMiddleware(secret string) *authMiddleware {
	return &authMiddleware{
		jwtSecret: []byte(secret),
	}
}

func (a *authMiddleware) ValidateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			// Err will trigger if the cookie doesn't exist
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized Request: Missing cookie"))
			return
		}

		// Delegate the complex parsing logic to our utility
		userID, email, err := helper.VerifyToken(tokenString, a.jwtSecret)
		if err != nil {
			// Cookie exists but token is tampered with or expired
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Session expired or invalid. Please log in again."))
			return
		}

		// Set the extracted variables into the Gin context
		c.Set("user_id", userID)
		c.Set("email", email)

		c.Next()
	}
}
