package main

import (
	"context"
	"fmt"
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

	ids, err := mong.InsertImageSetsMongo(ctx, database.Collection("imageset"), imageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = toxl.WriteToXL("ImageSets", ids)
	if err != nil {
		fmt.Println(err)
		return
	}

	allImageSets, imageSetMap, err := mong.GetImageSetMaps(ctx, database.Collection("imageset"))
	if err != nil {
		fmt.Println(err)
		return
	}

	exercises, err := fromxl.GetExercises(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	exerIDs, err := mong.InsertExerciseMongo(ctx, database.Collection("exercise"), exercises)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = toxl.WriteToXLDblRow("Exercises", exerIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	dynamics, err := fromxl.GetDynamics(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	dynamicIDs, err := mong.InsertDynamicMongo(ctx, database.Collection("dynamicstretch"), dynamics)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = toxl.WriteToXLDblRow("DynamicStretches", dynamicIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	statics, err := fromxl.GetStatics(allImageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	staticIDs, err := mong.InsertStaticMongo(ctx, database.Collection("staticstretch"), statics)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = toxl.WriteToXL("StaticStretches", staticIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	samples, err := fromxl.GetSampleAndProcess(database, ctx, imageSetMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	sampleIDs, err := mong.InsertSamples(ctx, database.Collection("sample"), samples)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = toxl.WriteToXL("Samples", sampleIDs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Complete")

}
