package bot

import (
	baseModels "bot-engine/models/mongo"
)

// Bot represents the core business entity for a Telegram Bot.
type Bot struct {
	baseModels.BaseModel `bson:",inline"`

	OwnerID        string
	TokenEncrypted string
	TelegramBotID  int64
	Username       string
	Status         BotStatus
	WebhookURL     string
}

type BotStatus string

const (
	StatusActive   BotStatus = "active"
	StatusInactive BotStatus = "inactive"
	StatusRevoked  BotStatus = "revoked"
)
