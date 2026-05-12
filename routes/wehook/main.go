package wehook

import (
	webhook "bot-engine/handlers/engine"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, webHandler *webhook.WebhookHandler) {

	// Exposed endpoint that Telegram targets
	// E.g., POST https://api.yourdomain.com/webhooks/766dca0b-9ee8-4c12-9844-323e20ec4cf7
	router.POST("/:webhook_token", webHandler.HandleTelegramWebhook)
}
