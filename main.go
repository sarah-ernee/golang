package main

import (
	"context"
	"fmt"
	"log"
	"module/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func connectMongoDB(collectionName string) *mongo.Collection {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	uri := "mongodb+srv://dev:sgtunnel2024@dev-gcp-tunnel.u0j98ms.mongodb.net/?retryWrites=true&w=majority&appName=dev-gcp-tunnel"
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var connectErr error
	client, connectErr = mongo.Connect(context.TODO(), opts)
	if connectErr != nil {
		log.Panic("Failed to connect to client:", connectErr)
	}

	// Will watch the non-time series dedicated collection
	// Dedicated collection will only ever have one document?
	// Constantly being replaced when ring number +1 or insert - can do document level change
	collection := client.Database("sample_supplies").Collection(collectionName)
	return collection
}

// Function defaults to listening to document and field level changes
func SetupChangeStream(collectionName *mongo.Collection) (*mongo.ChangeStream, error) {
	// Change stream output is typically delta fields, data that is changed
	// Pipeline can be tweaked to ignore some change events (only INSERT eg) or fields returned
	pipeline := mongo.Pipeline{
		// Disregard "delete" events from triggering a change event
		// bson.D{
		// 	{Key: "$match", Value: bson.D{
		// 		{Key: "operationType", Value: bson.D{
		// 			{Key: "$in", Value: bson.A{"insert", "update", "replaced"}},
		// 		}},
		// 	}},
		// },

		// Listen for only document insertion event
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "operationType", Value: "insert"},
			}},
		},

		// Trigger only on specific field change
		// bson.D{
		// 	{Key: "$match", Value: bson.D{
		// 		{Key: "operationType", Value: "update"},
		// 		{Key: "updateDescription.updatedFields.storeLocation", Value: bson.D{
		// 			{Key: "$exists", Value: true},
		// 		}},
		// 	}},
		// },
		// bson.D{
		// 	{Key: "$project", Value: bson.D{
		// 		{Key: "fullDocument.embeddings", Value: 0},
		// 	}},
		// },

		// bson.D{
		// 	{Key: "$match", Value: bson.D{
		// 		{Key: "operationType", Value: "update"},
		// 	}},
		// },
	}
	streamOpts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	return collectionName.Watch(context.TODO(), pipeline, streamOpts)
}

func ListenForChanges(changeStream *mongo.ChangeStream) {
	for changeStream.Next(context.TODO()) {
		var changeEvent bson.M
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Println("Error decoding change event: ", err)
			continue
		}

		collection1 := connectMongoDB("test")
		helper.DummyTest(changeEvent, collection1)

		if updateDesc, ok := changeEvent["updateDescription"].(bson.M); ok {
			updatedFields := updateDesc["updatedFields"].(bson.M)
			fmt.Printf("Updated fields: %+v\n", updatedFields)
			// fmt.Printf("Updated fields: %+v\n", updateDesc["updatedFields"])

			// // Capturing previous and current value of updated field
			// if prevValues, ok := changeEvent["fullDocumentBeforeChange"].(bson.M); ok {
			// 	if prevStoreLocation, exists := prevValues["storeLocation"]; exists {
			// 		fmt.Printf("Previous storeLocation: %v\n", prevStoreLocation)
			// 	}
			// }
			// if newStoreLocation, exists := updatedFields["storeLocation"]; exists {
			// 	fmt.Printf("New storeLocation: %v\n", newStoreLocation)
			// }
		}

		// fmt.Println("--------------------")
		// fmt.Printf("Operation Type: %v\n", changeEvent["operationType"])
		// if documentKey, ok := changeEvent["documentKey"].(bson.M); ok {
		// 	fmt.Printf("ID of document changed: %v\n", documentKey["_id"])
		// }

		// // Print out fields of changed document
		// if fullDoc, ok := changeEvent["fullDocument"].(bson.M); ok {
		// 	for key, value := range fullDoc {
		// 		fmt.Printf("%s: %v\n ", key, value)
		// 	}
		// } else {
		// 	log.Println("Document not found in change event")
		// }
	}

	if err := changeStream.Err(); err != nil {
		log.Println("Failure occured during change stream: ", err)
	}
}

func main() {
	collection2 := connectMongoDB("sales")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Failed to disconnect client: ", err)
		}
	}()

	changeStream, err := SetupChangeStream(collection2)
	if err != nil {
		log.Panic("Error occured before change stream: ", err)
	}

	// Stream closes when cursor is explicitly closed with timeout
	// Or invalidate event aka collection cannot be found or does not exist
	defer changeStream.Close(context.TODO())
	ListenForChanges(changeStream)
}
