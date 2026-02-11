package configs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(EnvMongoURI())

	// Configure TLS
	clientOptions.SetTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	// Set timeouts
	clientOptions.SetServerSelectionTimeout(30 * time.Second)
	clientOptions.SetConnectTimeout(30 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)

	// Configure connection pool
	clientOptions.SetMaxPoolSize(100)                  // Maximum number of connections in the pool
	clientOptions.SetMinPoolSize(5)                    // Minimum number of connections in the pool
	clientOptions.SetMaxConnIdleTime(30 * time.Minute) // Maximum time a connection can be idle before being closed

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(EnvMongoDB()).Collection(collectionName)
	return collection
}
