package users

import (
	userHandler "di/handlers/users"

	"github.com/gin-gonic/gin"
)

func SetupUserRouters(router *gin.RouterGroup, handler *userHandler.UserHandler) {

	router.GET("/", handler.List)

	router.POST("/", handler.Create)

}
