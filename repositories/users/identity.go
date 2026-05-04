package users

import (
	"context"
	"errors"
	"log"
	"time"

	userModels "bot-engine/models/mongo/users"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// The Contract
type IdentityRepository interface {
	Create(ctx context.Context, data *userModels.Identity) error
	GetByProvider(ctx context.Context, provider, providerUserID string) (*userModels.Identity, error)
}

type identityRepository struct {
	Collection *mongo.Collection
}

// Constructor
func NewIdentityRepository(db *mongo.Database) IdentityRepository {
	repo := &identityRepository{
		Collection: db.Collection("identities"),
	}

	repo.ensureIndexes()
	return repo
}

// ensureIndexes prevents the same Telegram account from being linked to multiple users
func (r *identityRepository) ensureIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO
	// Create a compound unique index on Provider + ProviderUserID
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "provider", Value: 1},
			{Key: "provider_user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("WARNING: Failed to create unique index for identities: %v\n", err)
	}
}

// Create inserts a new linked account
func (r *identityRepository) Create(ctx context.Context, data *userModels.Identity) error {
	data.ID = bson.NilObjectID
	data.LinkedAt = time.Now()

	created, err := r.Collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	if id, ok := created.InsertedID.(bson.ObjectID); ok {
		data.ID = id
	}

	return nil
}

// GetByProvider finds if a specific third-party account is already linked
func (r *identityRepository) GetByProvider(ctx context.Context, provider, providerUserID string) (*userModels.Identity, error) {
	var identity userModels.Identity

	filter := bson.D{
		{Key: "provider", Value: provider},
		{Key: "provider_user_id", Value: providerUserID},
	}

	err := r.Collection.FindOne(ctx, filter).Decode(&identity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Return nil so the service knows it's a new account
		}
		return nil, err
	}

	return &identity, nil
}
