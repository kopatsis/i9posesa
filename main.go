package main

import (
	"context"
	"fmt"
	"i9posesa/assets"
	"i9posesa/fromxl"
	"i9posesa/mong"
	"i9posesa/toxl"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx := context.Background()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("i9pos")

	fmt.Println("ImageSets: ")
	imageSets, err := fromxl.GetImageSets()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ImageSets Write Mongo: ")
	ids, err := mong.InsertImageSetsMongo(ctx, database.Collection("imageset"), imageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ImageSets Write XL: ")
	err = toxl.WriteToXL("ImageSets", ids)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ImageSets Gets: ")
	allImageSets, imageSetMap, err := mong.GetImageSetMaps(ctx, database.Collection("imageset"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Exers: ")
	exercises, err := fromxl.GetExercises(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Exers Write Mongo: ")
	exerIDs, err := mong.InsertExerciseMongo(ctx, database.Collection("exercise"), exercises)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Exers Write XL: ")
	err = toxl.WriteToXLDblRow("Exercises", exerIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Dynamics: ")
	dynamics, err := fromxl.GetDynamics(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Dynamics Write Mongo: ")
	dynamicIDs, err := mong.InsertDynamicMongo(ctx, database.Collection("dynamicstretch"), dynamics)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Dynamics Write XL: ")
	err = toxl.WriteToXLDblRow("DynamicStretches", dynamicIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Statics: ")
	statics, err := fromxl.GetStatics(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Statics Write Mongo: ")
	staticIDs, err := mong.InsertStaticMongo(ctx, database.Collection("staticstretch"), statics)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Statics Write XL: ")
	err = toxl.WriteToXL("StaticStretches", staticIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Samples: ")
	samples, err := fromxl.GetSampleAndProcess(database, ctx, imageSetMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Samples Write Mongo: ")
	sampleIDs, err := mong.InsertSamples(ctx, database.Collection("sample"), samples)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Samples Write XL: ")
	err = toxl.WriteToXL("Samples", sampleIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Transitions: ")
	transition, err := fromxl.GetTransitions(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Transitions Write Mongo: ")
	transitionIDs, err := mong.InsertTransition(ctx, database.Collection("transition"), []assets.TransitionMatrix{transition})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Transitions Write XL: ")
	err = toxl.WriteTransitionXL(transitionIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Complete")

}
