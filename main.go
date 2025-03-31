package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Sample struct {
	ID   bson.ObjectID `bson:"_id,omitempty"`
	Name string        `bson:"name"`
	Age  int           `bson:"age"`
}

func main() {
	ctx := context.Background()

	////////////////////////
	// Connect to MongoDB //
	////////////////////////
	client, err := mongo.Connect(options.Client().SetAuth(options.Credential{
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("firstdb").Collection("firstcollection")

	/////////////////////
	// Insert Document //
	/////////////////////

	r, err := collection.InsertOne(ctx, Sample{Name: "John", Age: 30})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted ID: %v\n", r.InsertedID)

	///////////////////
	// Find Document //
	///////////////////

	filter := bson.M{"name": "John"}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	var results []Sample
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	fmt.Printf("[BEFORE UPDATE] Found %d documents\n", len(results))
	for _, r := range results {
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", r.ID, r.Name, r.Age)
	}

	/////////////////////
	// Update Document //
	/////////////////////

	results[0].Name = "Jane"

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": results[0]})
	if err != nil {
		panic(err)
	}

	cursor2, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	var results2 []Sample
	if err = cursor2.All(ctx, &results2); err != nil {
		panic(err)
	}

	fmt.Printf("[AFTER UPDATE] Found %d documents\n", len(results2))
	for _, r := range results2 {
		fmt.Printf("ID: %s, Name: %s, Age: %d\n", r.ID, r.Name, r.Age)
	}

	//////////////////////
	// Delete Documents //
	//////////////////////

	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted document(s)")
}
