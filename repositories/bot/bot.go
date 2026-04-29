package bot

import (
	botModels "bot-engine/models/mongo/bot"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Bot = botModels.Bot

type BotRepository interface {
	Create(ctx context.Context, bot *Bot) error
	GetByOwnerID(ctx context.Context, ownerID string) ([]Bot, error)
}

type botRepository struct {
	Collection *mongo.Collection
}

func NewBotRepository(db *mongo.Database) *botRepository {

	col := db.Collection("bots")

	return &botRepository{
		Collection: col,
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
