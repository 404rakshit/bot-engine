package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	ValidateRequest() gin.HandlerFunc
}

type authMiddleware struct {
}

func NewAuthMiddleware() *authMiddleware {
	return &authMiddleware{}
}

func (a *authMiddleware) ValidateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {

		// token := c.GetHeader("Authorization")

		// if token == "" {
		// 	c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized Request"))
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
