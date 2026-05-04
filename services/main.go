package services

import (
	"bot-engine/config"
	"bot-engine/repositories"
	"bot-engine/services/auth"
	"bot-engine/services/bot"
	"bot-engine/services/encryption"
	"bot-engine/services/telegram"
	"bot-engine/services/users"
)

type Services struct {
	UserService       users.UserService
	BotService        bot.BotService
	EncryptionService encryption.EncryptionService
	TelegramService   telegram.TelegramService
	AuthService       auth.AuthService
}

func NewServices(
	r *repositories.Repositories,
	botCfg *config.BotSecretsConfig,
	authCfg *config.AuthSecretsConfig,
) *Services {
	// 1. Instantiate the standalone services first (the ones with no service dependencies)
	// You might need to pass environment variables here later (like your AES key)
	encService := encryption.NewEncryptionService(botCfg)
	tgService := telegram.NewTelegramService()

	authSvc := auth.NewAuthService(
		r.UserRepository,
		r.IdentityRepository,
		authCfg.JWTSecret,
	)

	// 2. Instantiate the bot service, injecting the repository AND the services we just created
	bService := bot.NewBotService(
		r.BotRepository,
		tgService,  // Injected instance
		encService, // Injected instance
	)

	// 3. Instantiate the user service
	uService := users.NewUserService(r.UserRepository)

	// 4. Return the aggregated struct containing all fully-wired instances
	return &Services{
		UserService:       uService,
		BotService:        bService,
		EncryptionService: encService,
		TelegramService:   tgService,
		AuthService:       authSvc,
	}
}
