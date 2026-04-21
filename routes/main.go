package routes

import (
	"bot-engine/handlers"
	middleware "bot-engine/middlewares"
	"bot-engine/utils"
	"net/http"

	userRoutes "bot-engine/routes/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handlers.Handlers, m *middleware.Middlerware) *gin.Engine {
	router := gin.Default()
	config := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000/",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}

	router.Use(cors.New(config))

	router.GET("/", func(c *gin.Context) {
		response := utils.SuccessResponse(nil, "Server is Alive...")
		c.JSON(http.StatusOK, response)
	})

	apiGroup := router.Group("/api")

	apiGroup.Use(m.AuthMiddleware.ValidateRequest())

	{
		userGroup := apiGroup.Group("/users")
		userRoutes.SetupUserRouters(userGroup, h.UserHandler)
	}

	router.NoRoute(ErrorHandler) // Changed to NoRoute for better 404 handling

	return router
}

func ErrorHandler(c *gin.Context) {
	path := c.Request.URL.Path
	message := "Path '" + path + "' does not exists"

	response := utils.ErrorResponse(message)
	c.JSON(http.StatusNotFound, response)

}
