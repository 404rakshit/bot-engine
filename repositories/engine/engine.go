package engine

import (
	"context"
	"errors"
	"time"

	engineModels "bot-engine/models/mongo/engine"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EngineRepository interface {
	GetWorkflow(ctx context.Context, workflowID bson.ObjectID) (*engineModels.Workflow, error)
	GetSession(ctx context.Context, botID bson.ObjectID, chatID int64) (*engineModels.UserSession, error)
	SaveSession(ctx context.Context, session *engineModels.UserSession) error
}

type engineRepository struct {
	workflows *mongo.Collection
	sessions  *mongo.Collection
}

func NewEngineRepository(db *mongo.Database) EngineRepository {
	repo := &engineRepository{
		workflows: db.Collection("workflows"),
		sessions:  db.Collection("user_sessions"),
	}
	repo.ensureIndexes()
	return repo
}

func (r *engineRepository) ensureIndexes() {
	// Create unique compound index so a user only has one session per active bot
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "bot_id", Value: 1},
			{Key: "chat_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	r.sessions.Indexes().CreateOne(ctx, indexModel)
}

func (r *engineRepository) GetWorkflow(ctx context.Context, workflowID bson.ObjectID) (*engineModels.Workflow, error) {
	var wf engineModels.Workflow
	err := r.workflows.FindOne(ctx, bson.M{"_id": workflowID}).Decode(&wf)
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

func (r *engineRepository) GetSession(ctx context.Context, botID bson.ObjectID, chatID int64) (*engineModels.UserSession, error) {
	var sess engineModels.UserSession
	filter := bson.M{"bot_id": botID, "chat_id": chatID}

	err := r.sessions.FindOne(ctx, filter).Decode(&sess)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // First contact: Session doesn't exist yet
		}
		return nil, err
	}
	return &sess, nil
}

func (r *engineRepository) SaveSession(ctx context.Context, session *engineModels.UserSession) error {
	session.UpdatedAt = time.Now().UTC()

	filter := bson.M{"bot_id": session.BotID, "chat_id": session.ChatID}
	updateOptions := options.Replace().SetUpsert(true)

	_, err := r.sessions.ReplaceOne(ctx, filter, session, updateOptions)
	return err
}
