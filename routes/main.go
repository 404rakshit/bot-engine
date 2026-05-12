package routes

import (
	"bot-engine/handlers"
	middleware "bot-engine/middlewares"
	"bot-engine/utils"
	"fmt"
	"net/http"
	"time"

	authRoutes "bot-engine/routes/auth"
	userRoutes "bot-engine/routes/users"
	webhookRoutes "bot-engine/routes/wehook"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Origin Header:", c.GetHeader("Origin"))
		c.Next()
	}
}

func NewRouter(h *handlers.Handlers, m *middleware.Middlewares) *gin.Engine {
	router := gin.Default()
	config := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://botapi.expdev.me",
			"https://share.expdev.me",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// router.Use(Logger())

	router.Use(cors.New(config))

	router.GET("/", func(c *gin.Context) {
		response := utils.SuccessResponse(nil, "Server is Alive...")
		c.JSON(http.StatusOK, response)
	})

	apiGroup := router.Group("/v1")

	{
		webhookGroup := apiGroup.Group("/webhooks")
		webhookRoutes.RegisterRoutes(webhookGroup, h.WebhookHandler)
	}

	{
		authGroup := apiGroup.Group("/auth")
		authRoutes.SetupAuthRouters(authGroup, h.AuthHandler)
	}

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
