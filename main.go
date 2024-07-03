package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := "mongodb+srv://dev:sgtunnel2024@dev-gcp-tunnel.u0j98ms.mongodb.net/?retryWrites=true&w=majority&appName=dev-gcp-tunnel"
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Panic("Failed to connect to client:", err)
	}
	database := client.Database("library")

	// Will watch the non-time series dedicated collection
	// Dedicated collection will only ever have one document?
	// Constantly being replaced when ring number +1
	collection = database.Collection("books")

}

func SetupChangeStream() (*mongo.ChangeStream, error) {
	// Change stream output is typically delta fields, data that is changed
	// Tweak pipeline to control the stream output
	pipeline := mongo.Pipeline{}
	streamOpts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	return collection.Watch(context.TODO(), pipeline, streamOpts)
}

// Function triggers for both document and field level changes
func ListenForChanges(changeStream *mongo.ChangeStream) {
	for changeStream.Next(context.TODO()) {
		var changeEvent bson.M
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Println("Error decoding change event: ", err)
			continue
		}

		fmt.Println("--------------------")
		fmt.Printf("Operation Type: %v\n", changeEvent["operationType"])

		switch changeEvent["operationType"] {
		case "insert", "update", "replace":
			if fullDoc, ok := changeEvent["fullDocument"].(bson.M); ok {
				fmt.Println("Change detected:")
				fmt.Println("Title: ", fullDoc["title"])
				fmt.Println("Author: ", fullDoc["authors"])
				fmt.Println("Year: ", fullDoc["year"])
				fmt.Println("Pages: ", fullDoc["pages"])
			} else {
				log.Println("Document not found in change event")
			}

		case "delete":
			if documentKey, ok := changeEvent["documentKey"].(bson.M); ok {
				fmt.Println("Document deleted:")
				fmt.Printf("Document ID: %v\n", documentKey["_id"])
			} else {
				log.Println("Document key not found in delete event")
			}

		default:
			fmt.Printf("Unhandled operation type: %v\n", changeEvent["operationType"])
		}
	}

	if err := changeStream.Err(); err != nil {
		log.Println("Failure occured during change stream: ", err)
	}
}

func main() {
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Failed to disconnect client: ", err)
		}
	}()

	changeStream, err := SetupChangeStream()
	if err != nil {
		log.Panic("Error occured before change stream: ", err)
	}

	// Stream closes when cursor is explicitly closed with timeout
	// Or invalidate event aka collection cannot be found or does not exist
	defer changeStream.Close(context.TODO())
	ListenForChanges(changeStream)
}
