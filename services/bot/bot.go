package bot

import (
	botModels "bot-engine/models/mongo/bot"
	botRepos "bot-engine/repositories/bot"
	encryptionService "bot-engine/services/encryption"
	"bot-engine/services/telegram"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EncryptionService = encryptionService.EncryptionService
type TelegramService = telegram.TelegramService

// Aliasing for cleaner code, matching your style
type botRepo = botRepos.BotRepository
type Bot = botModels.Bot

// BotService defines the core business capabilities for bots
type BotService interface {
	GetBotsByOwnerID(ctx context.Context, ownerID string) ([]Bot, error)
	ConnectNewBot(ctx context.Context, ownerID string, rawToken string) (*Bot, error)
}

type botService struct {
	botRepo           botRepo
	telegramService   TelegramService
	encryptionService EncryptionService
}

// NewBotService injects the database repository and external services (used by Wire)
func NewBotService(br botRepo, ts TelegramService, es EncryptionService) BotService {
	return &botService{
		botRepo:           br,
		telegramService:   ts,
		encryptionService: es,
	}
}

// GetBotsByOwnerID fetches all bots belonging to a specific user
func (s *botService) GetBotsByOwnerID(ctx context.Context, ownerID string) ([]Bot, error) {
	if ownerID == "" {
		return nil, errors.New("owner ID cannot be empty")
	}

	bots, err := s.botRepo.GetByOwnerID(ctx, ownerID)
	return bots, err
}

// ConnectNewBot orchestrates the secure bot connection business logic
func (s *botService) ConnectNewBot(ctx context.Context, ownerID string, rawToken string) (*Bot, error) {

	if rawToken == "" {
		return nil, errors.New("telegram token cannot be empty")
	}

	username, telegramBotID, err := s.telegramService.VerifyToken(ctx, rawToken)
	if err != nil {
		return nil, errors.New("failed to verify telegram token: " + err.Error())
	}

	encryptedToken, err := s.encryptionService.Encrypt(rawToken)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	ownerIDObj, err := bson.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, errors.New("id convertion issue: failed to convert into objectID")
	}

	// 3. Construct the Bot entity
	newBot := &Bot{
		OwnerID:        ownerIDObj,
		TokenEncrypted: encryptedToken,
		TelegramBotID:  telegramBotID,
		Username:       username,
		Status:         "active",
	}

	// 4. Optionally validate the model itself (if your model has a Validate method like User)
	// if err := newBot.Validate(); err != nil {
	// 	return nil, err
	// }

	// 5. Save to the database
	if err := s.botRepo.Create(ctx, newBot); err != nil {
		return nil, err
	}

	// 6. Return the created bot (Ensure your handler strips the encrypted token before sending JSON)
	return newBot, nil
}
