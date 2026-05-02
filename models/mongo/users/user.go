package users

import (
	baseModels "bot-engine/models/mongo"
	"errors"
	"strings"
)

type User struct {
	baseModels.BaseModel `bson:",inline"`

	Name         string `bson:"name" json:"name"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"password_hash" json:"-"`
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
