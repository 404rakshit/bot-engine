package main

import (
	"bot-engine/clients"
	"log"
)

func main() {
	log.Println("Just Print")

	client, err := clients.GetMongoClient()

	if err != nil {
		log.Fatal(err)
	}

	app, _ := InitializeAPI(client)
	app.Run(":8080")

}
