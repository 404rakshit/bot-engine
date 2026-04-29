package bot

import (
	botModels "bot-engine/models/mongo/bot"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Bot = botModels.Bot

type BotRepository interface {
	Create(ctx context.Context, bot *Bot) error
	GetByOwnerID(ctx context.Context, ownerID string) ([]Bot, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*Bot, error)
}

type botRepository struct {
	Collection *mongo.Collection
}

func NewBotRepository(db *mongo.Database) *botRepository {

	col := db.Collection("bots")

	repo := &botRepository{
		Collection: col,
	}

	repo.ensureIndexes()

	return repo
}

// ensureIndexes builds the required MongoDB indexes for the bot collection
func (r *botRepository) ensureIndexes() {
	// Give it a 10-second timeout so it doesn't block startup forever if the DB is slow
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Define the Index Model
	indexModel := mongo.IndexModel{
		// Set the field to index, and "1" for ascending order
		Keys: bson.D{{Key: "telegram_bot_id", Value: 1}},

		// 2. Enforce the Unique constraint
		Options: options.Index().SetUnique(true),
	}

	// 3. Execute the creation
	_, err := r.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// We use log.Printf instead of panic here. If an index fails to build,
		// you usually want the app to still boot so you can investigate,
		// but you want a loud warning in your logs.
		log.Printf("WARNING: Failed to create unique index for telegram_bot_id: %v\n", err)
	} else {
		log.Println("Successfully verified unique index on telegram_bot_id")
	}
}

func (r *botRepository) GetByOwnerID(ctx context.Context, ownerID string) ([]Bot, error) {
	var bots []Bot

	objID, err := bson.ObjectIDFromHex(ownerID)

	if err != nil {
		// Return an error if the provided string is not a valid 24-character hex string
		return nil, fmt.Errorf("invalid owner ID format: %w", err)
	}

	// 2. Use the primitive.ObjectID in the filter
	filter := bson.D{{Key: "owner_id", Value: objID}}

	// Find all documents
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed when the function returns
	defer cursor.Close(ctx)

	// Decode all documents into the bot slice
	if err := cursor.All(ctx, &bots); err != nil {
		return nil, err
	}

	// FIX: Return the actual slice of bots, and nil for the error!
	return bots, nil
}

func (r *botRepository) Create(ctx context.Context, data *Bot) error {

	data.ID = bson.NilObjectID
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	created, err := r.Collection.InsertOne(ctx, data)

	if err != nil {
		return err
	}

	ID, _ := created.InsertedID.(bson.ObjectID)

	data.ID = ID

	return nil
}

func (r *botRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*Bot, error) {
	var bot Bot

	filter := bson.D{{Key: "telegram_bot_id", Value: telegramID}}

	err := r.Collection.FindOne(ctx, filter).Decode(&bot)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &bot, nil
}
