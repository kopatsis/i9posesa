package mong

import (
	"context"
	"fmt"
	"i9posesa/assets"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertSamples(ctx context.Context, database *mongo.Database, structs []assets.Sample) ([]primitive.ObjectID, error) {

	collection := database.Collection("sample")

	ids := []primitive.ObjectID{}

	for _, item := range structs {
		var id primitive.ObjectID
		if item.ID.IsZero() {
			// Insert the item
			res, err := collection.InsertOne(ctx, item)
			if err != nil {
				fmt.Println(item.Name)
				return nil, err
			}

			id = res.InsertedID.(primitive.ObjectID)
			item.ID = id
			if err := updateActualEntry(item, database); err != nil {
				fmt.Println(item.Name)
				return nil, err
			}
			ids = append(ids, id)
		} else {
			// Update the item
			filter := bson.M{"_id": item.ID}
			_, err := collection.ReplaceOne(ctx, filter, item)
			if err != nil {
				return nil, err
			}
			if err := updateActualEntry(item, database); err != nil {
				fmt.Println(item.Name)
				return nil, err
			}
		}

	}

	return ids, nil
}

func updateActualEntry(sample assets.Sample, database *mongo.Database) error {
	collection := database.Collection("dynamicstretch")
	if sample.Type == "Exercise" {
		collection = database.Collection("exercise")
	} else if sample.Type == "Static Stretch" {
		collection = database.Collection("staticstretch")
	}

	filter := bson.M{"name": sample.Name}

	update := bson.M{
		"$set": bson.M{"sampleid": sample.ID.Hex()},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
