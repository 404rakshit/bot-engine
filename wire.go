//go:build wireinject

package main

import (
	"di/db"
	"di/handlers"
	middleware "di/middlewares"
	"di/repositories"
	"di/routes"
	"di/services"

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
