package repositories

import (
	"bot-engine/repositories/bot"
	"bot-engine/repositories/users"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repositories struct {
	UserRepository     users.UserRepository
	IdentityRepository users.IdentityRepository
	BotRepository      bot.BotRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		UserRepository:     users.NewUserRepository(db),
		BotRepository:      bot.NewBotRepository(db),
		IdentityRepository: users.NewIdentityRepository(db),
	}
}
