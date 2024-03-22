package fromxl

import (
	"errors"
	"fmt"
	"i9posesa/assets"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetExercises(imageSetMap map[string]string) ([]assets.Exercise, error) {
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

	exers := []assets.Exercise{}

	for i := 3; i < 1000; i++ {
		name, err := f.GetCellValue("Exercises", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if name == "" {
			break
		}

		maxSecsStr, err := f.GetCellValue("Exercises", "C"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if maxSecsStr == "" {
			posList2, err := getPosList(strconv.Itoa(i), f, imageSetMap)
			if err != nil {
				return nil, err
			}

			exers[len(exers)-1].PositionSlice2 = posList2
		} else {
			maxSecs, err := strconv.ParseFloat(maxSecsStr, 32)
			if err != nil {
				return nil, err
			}

			minSecsStr, err := f.GetCellValue("Exercises", "D"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			minSecs, err := strconv.ParseFloat(minSecsStr, 32)
			if err != nil {
				return nil, err
			}

			parent, err := f.GetCellValue("Exercises", "E"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			imagesetinit, err := f.GetCellValue("Exercises", "E"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			posList1, err := getPosList(strconv.Itoa(i), f, imageSetMap)
			if err != nil {
				return nil, err
			}

			imageSetID0, ok := "", false
			if imagesetinit != "" {
				imageSetID0, ok = imageSetMap[imagesetinit]
				if !ok {
					return nil, errors.New("image set name not in existing list in db")
				}
			} else {
				imageSetID0 = posList1[0].ImageSetID
			}

			exers = append(exers, assets.Exercise{
				Name:           name,
				MaxSecs:        float32(maxSecs),
				MinSecs:        float32(minSecs),
				Parent:         parent,
				ImageSetID0:    imageSetID0,
				PositionSlice1: posList1,
			})

			id, err := f.GetCellValue("Exercises", "A"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			if id != "" {

				primID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return nil, err
				}
				exers[len(exers)-1].ID = primID
			}

		}

	}

	return exers, nil
}

func getPosList(row string, f *excelize.File, imageSetMap map[string]string) ([]assets.ExerPosition, error) {

	ret := []assets.ExerPosition{}

	startCol := 7

	columnName, err := excelize.ColumnNumberToName(startCol)
	if err != nil {
		return nil, err
	}

	hasMore, err := f.GetCellValue("Exercises", columnName+row)
	if err != nil {
		return nil, err
	}

	for hasMore != "" {

		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		posName, err := f.GetCellValue("Exercises", columnName+row)
		if err != nil {
			return nil, err
		}

		imageSetID, ok := imageSetMap[posName]
		if !ok {
			return nil, errors.New("image set name not in existing list in db")
		}

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		hardcoded, err := f.GetCellValue("Exercises", columnName+row)
		if err != nil {
			return nil, err
		}

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		hardcodedSecsStr, err := f.GetCellValue("Exercises", columnName+row)
		if err != nil {
			return nil, err
		}

		hardcodedSecs, err := strconv.ParseFloat(hardcodedSecsStr, 32)
		if err != nil {
			return nil, err
		}

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		maxSecsStr, err := f.GetCellValue("Exercises", columnName+row)
		if err != nil {
			return nil, err
		}

		maxSecs, err := strconv.ParseFloat(maxSecsStr, 32)
		if err != nil {
			return nil, err
		}

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		percentSecsStr, err := f.GetCellValue("Exercises", columnName+row)
		if err != nil {
			return nil, err
		}

		percentSecs, err := strconv.ParseFloat(percentSecsStr, 32)
		if err != nil {
			return nil, err
		}

		ret = append(ret, assets.ExerPosition{
			ImageSetID:    imageSetID,
			Hardcoded:     hardcoded != "",
			HardcodedSecs: float32(hardcodedSecs),
			MaxSecs:       float32(maxSecs),
			PercentSecs:   float32(percentSecs),
		})

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		hasMore, err = f.GetCellValue("Exercises", columnName+row)
		if err != nil {
			return nil, err
		}

	}

	return ret, nil
}
