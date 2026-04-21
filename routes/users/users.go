package users

import (
	userHandler "bot-engine/handlers/users"

	"github.com/gin-gonic/gin"
)

func SetupUserRouters(router *gin.RouterGroup, handler *userHandler.UserHandler) {

	router.GET("/", handler.List)

	router.POST("/", handler.Create)

}
