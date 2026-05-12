package engine

import (
	baseModels "bot-engine/models/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NodeType string

const (
	NodeTrigger   NodeType = "trigger"
	NodeMessage   NodeType = "message"
	NodeInput     NodeType = "input"
	NodeCondition NodeType = "condition"
)

type Node struct {
	baseModels.BaseModel `bson:",inline" json:",inline"`

	Type         NodeType          `bson:"type" json:"type"`
	Content      string            `bson:"content,omitempty" json:"content,omitempty"`
	VariableName string            `bson:"variable_name,omitempty" json:"variable_name,omitempty"`
	Expression   string            `bson:"expression,omitempty" json:"expression,omitempty"`
	Next         string            `bson:"next,omitempty" json:"next,omitempty"`
	Branches     map[string]string `bson:"branches,omitempty" json:"branches,omitempty"`
}

type Workflow struct {
	baseModels.BaseModel `bson:",inline" json:",inline"`

	BotID       bson.ObjectID   `bson:"bot_id" json:"bot_id"` // Links workflow to specific bot
	Version     int             `bson:"version" json:"version"`
	StartNodeID string          `bson:"start_node_id" json:"start_node_id"`
	Nodes       map[string]Node `bson:"nodes" json:"nodes"`
}
