package repositories

import (
	"bot-engine/repositories/users"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	UserRepository users.UserRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
