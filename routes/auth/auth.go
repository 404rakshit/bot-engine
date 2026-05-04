package auth

import (
	authHandler "bot-engine/handlers/auth"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouters(router *gin.RouterGroup, handler *authHandler.AuthHandler) {

	router.POST("/login", handler.LoginUser)
	router.POST("/signup", handler.SignupUser)

	router.GET("/logout", handler.LogoutUser)

	router.POST("/telegram/login", handler.TelegramLogin)

}
