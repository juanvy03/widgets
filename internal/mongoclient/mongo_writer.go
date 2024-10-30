package mongoclient

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoProcessRegistration(collection string, event map[string]interface{}) {
	ctx := context.Background()

	uri := "mongodb://root:password@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("Widgets").Collection(collection)
	result, err := coll.InsertOne(ctx, event)
	if err != nil {
		panic(err)
	}
	log.Println(result)

}

func MongoProcessDeregistration(event map[string]interface{}) {
	ctx := context.Background()

	uri := "mongodb://root:password@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Println(event["serial_number"])

	filter := bson.D{{"serial_number", event["serial_number"]}}
	update := bson.D{{"$set", event}}
	//widgetData :=

	coll := client.Database("Widgets").Collection("registration")
	result := coll.FindOneAndUpdate(ctx, filter, update)

	// Declare the document to decode into
	var updatedDocument bson.M

	// Decode the result
	err = result.Decode(&updatedDocument)
	if err != nil {
		log.Fatalf("Error decoding updated document: %v", err)
	}

	log.Println(result)

}

func  MongoProcessLink(event map[string]interface{}) {
	/* 
	
		check if ports are available for both widgets 
	
	*/
	// FindByFilter widget1 and widget2

	// check if ports requested are true or false
	// if so then link them
	// else fail 

}