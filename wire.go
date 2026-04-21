//go:build wireinject

package main

import (
	"bot-engine/db"
	"bot-engine/handlers"
	middleware "bot-engine/middlewares"
	"bot-engine/repositories"
	"bot-engine/routes"
	"bot-engine/services"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitializeAPI(client *mongo.Client) (*gin.Engine, error) {

	wire.Build(

		db.NewMongoDatabase,

		repositories.NewRepositories,

		services.NewServices,

		handlers.NewHandler,

		middleware.NewMiddlerware,

		routes.NewRouter,
	)

	return nil, nil
}
