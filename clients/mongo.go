package clients

import (
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	// clientInstance holds the actual connection
	clientInstance *mongo.Client
	// clientError holds any error that occurred during connection
	clientError error
	// mongoOnce ensures the connection logic runs only once
	mongoOnce sync.Once
)

// GetMongoClient returns the singleton instance of the MongoDB client.
// It initializes the connection lazily on the first call.
func GetMongoClient() (*mongo.Client, error) {
	// once.Do guarantees the function inside is executed only once
	mongoOnce.Do(func() {
		uri := os.Getenv("MONGODB_URI")

		if uri == "" {
			uri = "mongodb://localhost:27017"
		}

		log.Print(uri)
		if uri == "" {
			clientError = fmt.Errorf("set your 'MONGODB_URI' environment variable. See: usage-examples/#environment-variable")
			return
		}

		// Connect to MongoDB
		// Note: In a real app, you might want to add a context with timeout here
		client, err := mongo.Connect(options.Client().ApplyURI(uri))
		if err != nil {
			clientError = err
			return
		}

		// Assign to the global variable
		clientInstance = client
	})

	return clientInstance, clientError
}
