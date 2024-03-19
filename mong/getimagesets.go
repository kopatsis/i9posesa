package mong

import (
	"context"
	"i9posesa/assets"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetImageSetMaps(ctx context.Context, collection *mongo.Collection) (map[string]string, map[string]assets.ImageSet, error) {
	imageSetMap := map[string]string{}
	allImageSets := map[string]assets.ImageSet{}

	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, nil, err
	}
	defer cur.Close(context.Background())

	// Iterate through the cursor to get all documents
	for cur.Next(context.Background()) {
		var result assets.ImageSet
		err := cur.Decode(&result)
		if err != nil {
			return nil, nil, err
		}
		imageSetMap[result.Name] = result.ID.Hex()
		allImageSets[result.ID.Hex()] = result
	}

	if err := cur.Err(); err != nil {
		return nil, nil, err
	}

	return imageSetMap, allImageSets, nil
}
