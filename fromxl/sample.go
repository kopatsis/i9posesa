package fromxl

import (
	"context"
	"errors"
	"fmt"
	"i9posesa/assets"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSampleAndProcess(database *mongo.Database, ctx context.Context, imageSetMap map[string]assets.ImageSet) ([]assets.Sample, error) {
	samples := []assets.Sample{}

	f, err := excelize.OpenFile("assets/posesaxl.xlsx")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for i := 3; i < 1000; i++ {
		name, err := f.GetCellValue("Samples", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if name == "" {
			break
		}

		description, err := f.GetCellValue("Samples", "C"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		dtype, err := f.GetCellValue("Samples", "D"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		samples = append(samples, assets.Sample{
			Name:        name,
			Description: description,
			Type:        dtype,
		})

		id, err := f.GetCellValue("Samples", "A"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if id != "" {

			primID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, err
			}
			samples[len(samples)-1].ID = primID
		}

		secsStr, err := f.GetCellValue("Samples", "E"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		secs64, err := strconv.ParseFloat(secsStr, 32)
		if err != nil {
			return nil, err
		}
		secs := float32(secs64)

		if dtype == "Exercise" {
			exer, err := getCorrespondingExer(ctx, database.Collection("exercise"), name)
			if err != nil {
				return nil, err
			}
			samples[len(samples)-1].ExOrStID = exer.ID.Hex()
			samples[len(samples)-1].Rep = createRepExer(secs, exer, imageSetMap)

		} else if dtype == "StaticStretch" {
			static, err := getCorrespondingStatic(ctx, database.Collection("staticstretch"), name)
			if err != nil {
				return nil, err
			}
			samples[len(samples)-1].ExOrStID = static.ID.Hex()
			samples[len(samples)-1].Rep = createRepStatic(secs, static, imageSetMap)

		} else if dtype == "DynamicStretch" {
			dynamic, err := getCorrespondingDynamic(ctx, database.Collection("dynamicstretch"), name)
			if err != nil {
				return nil, err
			}
			samples[len(samples)-1].ExOrStID = dynamic.ID.Hex()
			samples[len(samples)-1].Rep = createRepDynamic(secs, dynamic, imageSetMap)

		} else {
			return nil, errors.New("provided type doesn't exist")
		}

	}

	return samples, nil
}

func createRepExer(secs float32, exer assets.Exercise, imageSetMap map[string]assets.ImageSet) assets.Rep {
	ret := assets.Rep{}

	if len(exer.PositionSlice2) == 0 {

		positions, times := [][]string{}, []float32{}
		copySecs := secs

		for _, pos := range exer.PositionSlice1 {
			if pos.Hardcoded {
				copySecs -= pos.HardcodedSecs
			}
		}
		for _, pos := range exer.PositionSlice1 {
			if pos.Hardcoded {
				times = append(times, pos.HardcodedSecs)
			} else {
				times = append(times, pos.PercentSecs*copySecs)
			}
			positions = append(positions, imageSetMap[pos.ImageSetID].Mid)
		}

		ret.Positions = positions
		ret.FullTime = secs
		ret.Times = times
	} else {

		positions, times := [][]string{}, []float32{}
		copySecs := secs

		for _, pos := range exer.PositionSlice1 {
			if pos.Hardcoded {
				copySecs -= pos.HardcodedSecs
			}
		}
		for _, pos := range exer.PositionSlice2 {
			if pos.Hardcoded {
				copySecs -= pos.HardcodedSecs
			}
		}

		for _, pos := range exer.PositionSlice1 {
			if pos.Hardcoded {
				times = append(times, pos.HardcodedSecs)
			} else {
				times = append(times, 0.5*pos.PercentSecs*copySecs)
			}
			positions = append(positions, imageSetMap[pos.ImageSetID].Mid)
		}
		for _, pos := range exer.PositionSlice2 {
			if pos.Hardcoded {
				times = append(times, pos.HardcodedSecs)
			} else {
				times = append(times, 0.5*pos.PercentSecs*copySecs)
			}
			positions = append(positions, imageSetMap[pos.ImageSetID].Mid)
		}

		ret.Positions = positions
		ret.FullTime = secs
		ret.Times = times
	}

	return ret
}

func createRepDynamic(secs float32, stretch assets.DynamicStr, imageSetMap map[string]assets.ImageSet) assets.Rep {
	ret := assets.Rep{}

	if len(stretch.PositionSlice2) == 0 {

		positions, times := [][]string{}, []float32{}
		for _, pos := range stretch.PositionSlice1 {
			times = append(times, pos.PercentSecs*secs)
			positions = append(positions, imageSetMap[pos.ImageSetID].Mid)
		}

		ret.Positions = positions
		ret.FullTime = secs
		ret.Times = times
	} else {

		positions, times := [][]string{}, []float32{}
		for _, pos := range stretch.PositionSlice1 {
			times = append(times, 0.5*pos.PercentSecs*secs)
			positions = append(positions, imageSetMap[pos.ImageSetID].Mid)
		}
		for _, pos := range stretch.PositionSlice2 {
			times = append(times, 0.5*pos.PercentSecs*secs)
			positions = append(positions, imageSetMap[pos.ImageSetID].Mid)
		}

		ret.Positions = positions
		ret.FullTime = secs
		ret.Times = times
	}

	return ret
}

func createRepStatic(secs float32, stretch assets.StaticStr, imageSetMap map[string]assets.ImageSet) assets.Rep {

	ret := assets.Rep{}

	if stretch.ImageSetID2 != "" {
		positions := imageSetMap[stretch.ImageSetID1].Mid
		ret.Positions = [][]string{positions}
		ret.FullTime = secs
		ret.Times = []float32{secs}
	} else {
		positions1 := imageSetMap[stretch.ImageSetID1].Mid
		positions2 := imageSetMap[stretch.ImageSetID2].Mid
		ret.Positions = [][]string{positions1, positions2}
		ret.FullTime = secs
		ret.Times = []float32{secs / 2, secs / 2}
	}

	return ret
}

func getCorrespondingExer(ctx context.Context, collection *mongo.Collection, name string) (assets.Exercise, error) {
	var result assets.Exercise
	err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return assets.Exercise{}, err
		} else {
			return assets.Exercise{}, err
		}
	} else {
		return result, nil
	}
}

func getCorrespondingDynamic(ctx context.Context, collection *mongo.Collection, name string) (assets.DynamicStr, error) {
	var result assets.DynamicStr
	err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return assets.DynamicStr{}, err
		} else {
			return assets.DynamicStr{}, err
		}
	} else {
		return result, nil
	}
}

func getCorrespondingStatic(ctx context.Context, collection *mongo.Collection, name string) (assets.StaticStr, error) {
	var result assets.StaticStr
	err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return assets.StaticStr{}, err
		} else {
			return assets.StaticStr{}, err
		}
	} else {
		return result, nil
	}
}
