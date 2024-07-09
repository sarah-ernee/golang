package helper

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DummyTest(changeEvent bson.M, collection *mongo.Collection) {
	if fullDoc, ok := changeEvent["fullDocument"].(bson.M); ok {
		if id, ok := fullDoc["_id"].(string); ok {
			title := fullDoc["title"]

			newDoc := bson.D{
				{Key: "bookID", Value: id},
				{Key: "name", Value: title},
			}

			_, err := collection.InsertOne(context.TODO(), newDoc)
			if err != nil {
				log.Println("Error inserting document: ", err)
			} else {
				fmt.Println("Helper function successfully triggered by stream!")
			}
		}
	}
}
