package users

import (
	baseModels "bot-engine/models/mongo"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	baseModels.BaseModel `bson:",inline"`

	Name         string `bson:"name" json:"name"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"-"`
}

type Identity struct {
	baseModels.BaseModel `bson:",inline"`

	UserID         bson.ObjectID `bson:"user_id" json:"user_id"`
	Provider       string        `bson:"provider" json:"provider"`                 // e.g., "telegram", "google"
	ProviderUserID string        `bson:"provider_user_id" json:"provider_user_id"` // e.g., "123456789"
	LinkedAt       time.Time     `bson:"linked_at" json:"linked_at"`
}

func (u *User) Validate() error {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.PasswordHash == "" {
		return errors.New("password hash is required")
	}
	return nil
}
