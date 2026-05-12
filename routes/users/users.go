package users

import (
	userHandler "bot-engine/handlers/users"

	"github.com/gin-gonic/gin"
)

func SetupUserRouters(router *gin.RouterGroup, handler *userHandler.UserHandler) {

	router.GET("", handler.List)

	router.POST("", handler.Create)

	botGroup := router.Group("/:user_id/bots")

	{
		botGroup.GET("", handler.ListUserBots)
		botGroup.POST("", handler.CreateUserBot)
	}

}
