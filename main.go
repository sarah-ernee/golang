package main

import (
	"context"
	"fmt"
	"log"
	"module/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FullDocument struct {
	SaleDate time.Time `bson:"saleDate"`
	Customer customer  `bson:"customer"`
	Items    []item    `bson:"items"`
}

type customer struct {
	Gender string `bson:"gender"`
	Age    int    `bson:"age"`
	Email  string `bson:"email"`
}

type item struct {
	Name  string  `bson:"name"`
	Price float64 `bson:"price"`
}

type RingSummary struct {
	Timestamp  time.Time `bson:"timestamp"`
	MiningInfo Mining    `bson:"mining_info"`
}

type Mining struct {
	RingNumber int     `bson:"ring"`
	Chainage   float64 `bson:"chainage_head"`
}

var client *mongo.Client

func connectMongoDB(dbName, collectionName string) *mongo.Collection {
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
	collection := client.Database(dbName).Collection(collectionName)
	return collection
}

func SetupChangeStream(collectionName *mongo.Collection) (*mongo.ChangeStream, error) {
	// Change stream pipeline can watch for specific field changes or operation types
	// Pipeline can also manipulate change stream output
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "operationType", Value: "insert"},
			}},
		},

		// EXAMPLE 1: Exclude listening for DELETE event
		// bson.D{
		// 	{Key: "$match", Value: bson.D{
		// 		{Key: "operationType", Value: bson.D{
		// 			{Key: "$in", Value: bson.A{"insert", "update", "replaced"}},
		// 		}},
		// 	}},
		// },

		// EXAMPLE 2: Specific field level change
		// bson.D{
		// 	{Key: "$match", Value: bson.D{
		// 		{Key: "operationType", Value: "update"},
		// 		{Key: "updateDescription.updatedFields.storeLocation", Value: bson.D{
		// 			{Key: "$exists", Value: true},
		// 		}},
		// 	}},
		// },

		// EXAMPLE 3: Exclude embeddings field from output
		// bson.D{
		// 	{Key: "$project", Value: bson.D{
		// 		{Key: "fullDocument.embeddings", Value: 0},
		// 	}},
		// },
	}
	streamOpts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	return collectionName.Watch(context.TODO(), pipeline, streamOpts)
}

func ListenForChanges(changeStream *mongo.ChangeStream) {
	rawDataCol := connectMongoDB("sample_tbm", "tbm_sg_raw_nested")

	for changeStream.Next(context.TODO()) {
		var changeEvent bson.M
		if err := changeStream.Decode(&changeEvent); err != nil {
			log.Println("Error decoding change event: ", err)
		}

		// "Data capture" stage
		var result RingSummary
		opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
		err := rawDataCol.FindOne(context.Background(), bson.D{}, opts).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Println("No documents found in the collection")
			} else {
				log.Println("Error fetching latest document:", err)
			}
			continue
		}

		fmt.Println(result)
		fmt.Printf("Timestamp: %v, Ring Number: %d, Chainage: %.2f\n",
			result.Timestamp, result.MiningInfo.RingNumber, result.MiningInfo.Chainage)

		// pass into helper for data pump
		// fmt.Println("--------------------")
		// if documentKey, ok := changeEvent["documentKey"].(bson.M); ok {
		// 	fmt.Printf("CHANGED DOCUMENT ID: %v\n", documentKey["_id"])
		// }
		// fmt.Println("--------------------")

		// var doc FullDocument
		// bsonBytes, err := bson.Marshal(changeEvent)
		// if err != nil {
		// 	fmt.Println("Error marshalling primitive.M to BSON: ", err)
		// }
		// if err := bson.Unmarshal(bsonBytes, &doc); err == nil {
		// 	fmt.Println(doc.Customer.Email)

		// 	if len(doc.Items) > 0 {
		// 		fmt.Println(doc.Items[0].Name)
		// 	} else {
		// 		fmt.Println("No items array")
		// 	}
		// }

		// "Data pump" stage
		collection1 := connectMongoDB("sample_supplies", "test")
		helper.DummyTest(changeEvent, collection1)
	}

	if err := changeStream.Err(); err != nil {
		log.Println("Failure occured during change stream: ", err)
	}
}

func main() {
	collection2 := connectMongoDB("sample_supplies", "sales")
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Failed to disconnect client: ", err)
		}
	}()

	changeStream, err := SetupChangeStream(collection2)
	if err != nil {
		log.Panic("Error occured before change stream: ", err)
	}
	defer changeStream.Close(context.TODO())
	ListenForChanges(changeStream)
}
