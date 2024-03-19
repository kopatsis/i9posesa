package mong

import (
	"context"
	"i9posesa/assets"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertTransition(ctx context.Context, collection *mongo.Collection, structs []assets.TransitionMatrix) ([]primitive.ObjectID, error) {

	ids := make([]primitive.ObjectID, len(structs))

	for i, item := range structs {
		var id primitive.ObjectID
		if item.ID.IsZero() {
			// Insert the item
			res, err := collection.InsertOne(ctx, item)
			if err != nil {
				return nil, err
			}
			id = res.InsertedID.(primitive.ObjectID)
		} else {
			// Update the item
			id = item.ID
			filter := bson.M{"_id": item.ID}
			_, err := collection.ReplaceOne(ctx, filter, item)
			if err != nil {
				return nil, err
			}
		}

		ids[i] = id
	}

	return ids, nil
}
