package auth

import (
	"di/utils"
	"net/http"

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

		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized Request"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// package auth

// import (
// 	"di/utils"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// // 1. Keep the Interface
// type AuthMiddleware interface {
// 	ValidateRequest() gin.HandlerFunc
// }

// // 2. Keep the concrete struct private
// type authMiddleware struct{}

// // 3. CHANGE: Return the Interface here
// // This tells Go: "Here is something that satisfies AuthMiddleware"
// func NewAuthMiddleware() AuthMiddleware {
// 	return &authMiddleware{}
// }

// func (a *authMiddleware) ValidateRequest() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := c.GetHeader("Authorization")
// 		if token == "" {
// 			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized Request"))
// 			c.Abort()
// 			return
// 		}
// 		c.Next()
// 	}
// }
