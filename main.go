package main

import (
	"context"
	"fmt"
	"i9posesa/assets"
	"i9posesa/fromxl"
	"i9posesa/mong"
	"i9posesa/toxl"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	ctx := context.Background()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectStr := os.Getenv("MONGOSTRING")
	clientOptions := options.Client().ApplyURI(connectStr).SetServerAPIOptions(serverAPI)

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
	imageSets, err := fromxl.GetImageSetNameMap()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Exers: ")
	exercises, err := fromxl.GetExercises(imageSets)
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
	dynamics, err := fromxl.GetDynamics(imageSets)
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
	statics, err := fromxl.GetStatics(imageSets)
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
	samples, err := fromxl.GetSampleAndProcess(database, ctx, imageSets)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Samples Write Mongo: ")
	sampleIDs, err := mong.InsertSamples(ctx, database, samples)
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
	transition, err := fromxl.GetTransitions(imageSets)
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
