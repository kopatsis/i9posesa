package mong

import (
	"context"
	"fmt"
	"i9posesa/assets"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertSamples(ctx context.Context, collection *mongo.Collection, structs []assets.Sample) ([]primitive.ObjectID, error) {

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
			ids = append(ids, id)
		} else {
			// Update the item
			filter := bson.M{"_id": item.ID}
			_, err := collection.ReplaceOne(ctx, filter, item)
			if err != nil {
				return nil, err
			}
		}

	}

	return ids, nil
}
