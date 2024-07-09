package helper

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DummyTest(changeEvent bson.M, collection *mongo.Collection) {
	newDoc := bson.D{
		{Key: "name", Value: "dummy"},
	}

	_, err := collection.InsertOne(context.TODO(), newDoc)
	if err != nil {
		log.Println("Error inserting document: ", err)
	} else {
		fmt.Println("Helper function successfully triggered by stream!")
	}

}
