package users

import (
	"context"
	userModels "di/models/mongo/users"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userModel = userModels.User

type UserRepository interface {
	List(ctx context.Context) ([]userModel, error)
	Create(ctx context.Context, data *userModel) error
}

type userRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {

	col := db.Collection("users")

	return &userRepository{
		Collection: col,
	}
}
func (r *userRepository) List(ctx context.Context) ([]userModel, error) {
	var users []userModel

	// Find all documents
	cursor, err := r.Collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed when the function returns
	defer cursor.Close(ctx)

	// Decode all documents into the users slice
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// FIX: Return the actual slice of users, and nil for the error!
	return users, nil
}

func (r *userRepository) Create(ctx context.Context, data *userModel) error {

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
