package config

import "bot-engine/utils"

type mongoUri = string

func GetMongoConfig() mongoUri {

	mongoUriCfg := utils.GetEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017")

	return mongoUriCfg
}
