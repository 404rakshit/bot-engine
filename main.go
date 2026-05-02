package main

import (
	"bot-engine/clients"
	"bot-engine/config"
	"bot-engine/utils"
	"log"
)

func main() {

	utils.LoadEnv(".env")

	mongoUri := config.GetMongoConfig()

	client, err := clients.GetMongoClient(mongoUri)

	if err != nil {
		log.Fatal(err)
	}

	app, _ := InitializeAPI(client)
	app.Run(":8000")

}
