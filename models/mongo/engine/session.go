package engine

import (
	baseModels "bot-engine/models/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserSession struct {
	baseModels.BaseModel `bson:",inline" json:",inline"`

	BotID         bson.ObjectID          `bson:"bot_id" json:"bot_id"`           // The bot context
	ChatID        int64                  `bson:"chat_id" json:"chat_id"`         // End-user Telegram ID
	WorkflowID    bson.ObjectID          `bson:"workflow_id" json:"workflow_id"` // Active workflow
	CurrentNodeID string                 `bson:"current_node_id" json:"current_node_id"`
	Context       map[string]interface{} `bson:"context" json:"context"` // Saved dynamic variables
}
