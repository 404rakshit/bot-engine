package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	botServices "bot-engine/services/bot"
	engineServices "bot-engine/services/engine"
	"bot-engine/utils"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	BotService    botServices.BotService       // Pure Service dependency
	EngineService engineServices.EngineService // Pure Service dependency
	HTTPClient    *http.Client
}

func NewWebhookHandler(
	botService botServices.BotService,
	engineService engineServices.EngineService,
) *WebhookHandler {
	return &WebhookHandler{
		BotService:    botService,
		EngineService: engineService,
		HTTPClient:    &http.Client{},
	}
}

// type TelegramWebhookHandler struct {
// 	EngineService engineServices.EngineService
// }

type TelegramUpdate struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func (h *WebhookHandler) HandleTelegramWebhook(c *gin.Context) {
	webhookToken := c.Param("webhook_token")
	if webhookToken == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Missing webhook token"))
		return
	}

	// 1. Ask the service layer to resolve and decrypt the bot context
	botDoc, rawBotToken, err := h.BotService.GetActiveBotByWebhook(c.Request.Context(), webhookToken)
	if err != nil {
		// Log error internally, but reply 200 OK so Telegram stops retrying
		c.JSON(http.StatusOK, gin.H{"status": "bot_resolution_failed"})
		return
	}

	// 2. Parse incoming Telegram Update
	var update TelegramUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	// 3. Process message through the engine
	replies, err := h.EngineService.ProcessIncomingMessage(
		c.Request.Context(),
		botDoc.ID,
		update.Message.Chat.ID,
		update.Message.Text,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "engine_failed"})
		return
	}

	// 4. Dispatch replies in background thread
	go h.dispatchReplies(rawBotToken, update.Message.Chat.ID, replies)

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func (h *WebhookHandler) dispatchReplies(token string, chatID int64, messages []string) {
	for _, text := range messages {
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
		payload := map[string]interface{}{
			"chat_id": chatID,
			"text":    text,
		}

		body, _ := json.Marshal(payload)
		resp, err := h.HTTPClient.Post(apiURL, "application/json", bytes.NewBuffer(body))
		if err == nil {
			resp.Body.Close()
		}
	}
}

// func (h *TelegramWebhookHandler) HandleWebhook(c *gin.Context) {
// 	// 1. Parse BotID context from route URL (e.g., /webhooks/:bot_id)
// 	botIDStr := c.Param("bot_id")
// 	botID, err := bson.ObjectIDFromHex(botIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid bot ID"))
// 		return
// 	}

// 	// 2. Parse Telegram Update
// 	var update TelegramUpdate
// 	if err := c.ShouldBindJSON(&update); err != nil {
// 		c.JSON(http.StatusOK, nil) // Respond 200 immediately so Telegram doesn't retry
// 		return
// 	}

// 	if update.Message.Text == "" {
// 		c.JSON(http.StatusOK, nil)
// 		return
// 	}

// 	// 3. Process message through our fast state machine engine
// 	messages, err := h.EngineService.ProcessIncomingMessage(
// 		c.Request.Context(),
// 		botID,
// 		update.Message.Chat.ID,
// 		update.Message.Text,
// 	)

// 	if err != nil {
// 		// Log internal engine issues but don't panic or block the endpoint
// 		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
// 		return
// 	}

// 	// 4. Send messages back to user (The Muscle Executor)
// 	// loop over 'messages' and fire HTTP calls to https://api.telegram.org/bot<Token>/sendMessage

// 	c.JSON(http.StatusOK, gin.H{"status": "processed", "replies_count": len(messages)})
// }
