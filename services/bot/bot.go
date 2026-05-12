package bot

import (
	botModels "bot-engine/models/mongo/bot"
	botRepos "bot-engine/repositories/bot"
	encryptionService "bot-engine/services/encryption"
	"bot-engine/services/telegram"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	GetActiveBotByWebhook(ctx context.Context, webhookToken string) (*botModels.Bot, string, error)
}

type botService struct {
	botRepo           botRepo
	telegramService   TelegramService
	encryptionService EncryptionService
	botWebhookService WebhookRegistrar
}

// NewBotService injects the database repository and external services (used by Wire)
func NewBotService(br botRepo, ts TelegramService, es EncryptionService, bw WebhookRegistrar) BotService {
	return &botService{
		botRepo:           br,
		telegramService:   ts,
		encryptionService: es,
		botWebhookService: bw,
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
		return nil, errors.New("Telegram token cannot be empty")
	}

	username, telegramBotID, err := s.telegramService.VerifyToken(ctx, rawToken)
	if err != nil {
		return nil, errors.New("Failed to verify telegram token: " + err.Error())
	}

	existingBot, err := s.botRepo.GetByTelegramID(ctx, telegramBotID)
	if err != nil {
		return nil, errors.New("Internal error checking for duplicate bots")
	}

	if existingBot != nil {
		// Bot exists! Reject the request.
		return nil, errors.New("This telegram bot is already connected to the system")
	}

	encryptedToken, err := s.encryptionService.Encrypt(rawToken)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	ownerIDObj, err := bson.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, errors.New("id convertion issue: failed to convert into objectID")
	}

	webhookToken := uuid.New().String()

	// 3. Construct the Bot entity
	newBot := &Bot{
		OwnerID:        ownerIDObj,
		TokenEncrypted: encryptedToken,
		TelegramBotID:  telegramBotID,
		Username:       username,
		WebhookToken:   webhookToken,
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

	// Note: Replace with your actual domain name!
	systemWebhookURL := fmt.Sprintf("https://bot.expdev.me/v1/webhooks/%s", webhookToken)

	if err := s.botWebhookService.Register(rawToken, systemWebhookURL); err != nil {
		return nil, fmt.Errorf("bot saved but failed to register with Telegram: %w", err)
	}
	// 6. Return the created bot (Ensure your handler strips the encrypted token before sending JSON)
	return newBot, nil
}

func (s *botService) GetActiveBotByWebhook(ctx context.Context, webhookToken string) (*botModels.Bot, string, error) {
	// 1. Fetch bot using repository
	botDoc, err := s.botRepo.GetByWebhookToken(ctx, webhookToken)
	if err != nil {
		return nil, "", fmt.Errorf("bot not found or inactive: %w", err)
	}

	// 2. Decrypt the bot token within the safety of the service layer
	rawToken, err := s.encryptionService.Decrypt(botDoc.TokenEncrypted)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decrypt bot token: %w", err)
	}

	return botDoc, rawToken, nil
}
