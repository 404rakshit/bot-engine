package mogno

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BaseModel struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
