package bot

import (
	baseModels "bot-engine/models/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Bot represents the core business entity for a Telegram Bot.
type Bot struct {
	baseModels.BaseModel `bson:",inline" json:",inline"`

	OwnerID        bson.ObjectID `bson:"owner_id" json:"owner_id,omitempty"`
	TokenEncrypted string        `bson:"token_encrypted" json:"token_encrypted"`
	TelegramBotID  int64         `bson:"telegram_bot_id" json:"telegram_bot_id"`
	Username       string        `bson:"username" json:"username"`
	Status         BotStatus     `bson:"status" json:"status"`
	WebhookURL     string        `bson:"webhook_url" json:"webhook_url"`
}

type BotStatus string

const (
	StatusActive   BotStatus = "active"
	StatusInactive BotStatus = "inactive"
	StatusRevoked  BotStatus = "revoked"
)
