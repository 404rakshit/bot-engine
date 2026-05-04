package auth

import (
	"bot-engine/helper"
	"bot-engine/utils"
	"fmt"
	"net/http"

	authModel "bot-engine/models/mongo/auth"

	"github.com/gin-gonic/gin"
)

type TelegramAuthPayload = authModel.TelegramAuthPayload

func (h *AuthHandler) TelegramLogin(c *gin.Context) {
	// 1. Bind to map for exact cryptographic validation
	var rawPayload map[string]interface{}
	if err := c.ShouldBindJSON(&rawPayload); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid payload format"))
		return
	}

	// 2. Convert everything to strings safely, avoiding scientific notation
	stringMap := make(map[string]string)
	for key, value := range rawPayload {
		switch v := value.(type) {
		case float64:
			// "%.0f" forces Go to print the float as a flat integer string
			// (e.g., 5932880028 instead of 5.932880028e+09)
			stringMap[key] = fmt.Sprintf("%.0f", v)
		default:
			// For standard strings (like first_name, hash, username)
			stringMap[key] = fmt.Sprintf("%v", v)
		}
	}

	// 2. Verify signature
	if err := helper.VerifyTelegramAuth(stringMap, h.BotCfg.BotToken); err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		return
	}

	// 3. Bind to strict struct for business logic
	var payload TelegramAuthPayload

	// Safely extract the ID and AuthDate since they arrived as float64 from JSON
	if id, ok := rawPayload["id"].(float64); ok {
		payload.ID = int64(id)
	}
	if authDate, ok := rawPayload["auth_date"].(float64); ok {
		payload.AuthDate = int64(authDate)
	}

	// We bind again using the c.Request.Body (ensure you handle reading the body twice
	// or just map the data manually from rawData to the struct)
	if fn, ok := rawPayload["first_name"].(string); ok {
		payload.FirstName = fn
	}
	if ln, ok := rawPayload["last_name"].(string); ok {
		payload.LastName = ln
	}
	if un, ok := rawPayload["username"].(string); ok {
		payload.Username = un
	}
	if photo, ok := rawPayload["photo_url"].(string); ok {
		payload.PhotoURL = photo
	}

	// 4. Process Login/Signup
	token, err, userIdStr := h.AuthService.ProcessTelegramLogin(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Authentication failed"))
		return
	}

	customResponse := CustomeUserResponse{
		UserID: userIdStr,
		Email:  "",
	}

	// 5. Set the HttpOnly Cookie just like your standard login
	setAuthCookie(c, token)
	c.JSON(http.StatusOK, utils.SuccessResponse(customResponse, "Telegram login successful"))
}
