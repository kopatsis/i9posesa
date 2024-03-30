package fromxl

import (
	"errors"
	"fmt"
	"i9posesa/assets"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTransitions(allImageSets map[string]string) (assets.TransitionMatrix, error) {
	matrix := assets.TransitionMatrix{}

	f, err := excelize.OpenFile("assets/posesaxl.xlsx")
	if err != nil {
		fmt.Println(err)
		return matrix, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	id, err := f.GetCellValue("Transitions", "B1")
	if err != nil {
		return matrix, err
	}

	if id != "" {
		primID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return matrix, err
		}
		matrix.ID = primID
	}

	for i := 3; i < 1000; i++ {

		row := strconv.Itoa(i)

		name, err := f.GetCellValue("Transitions", "D"+row)
		if err != nil {
			return matrix, err
		}

		if name == "" {
			break
		}

		slowRep := assets.TransitionRep{
			ImageSetIDs: []string{},
			Times:       []float32{},
			FullTime:    0,
		}
		regRep := assets.TransitionRep{
			ImageSetIDs: []string{},
			Times:       []float32{},
			FullTime:    0,
		}
		fastRep := assets.TransitionRep{
			ImageSetIDs: []string{},
			Times:       []float32{},
			FullTime:    0,
		}

		startCol := 4

		for name != "" {

			imageID, ok := allImageSets[name]
			if !ok {
				fmt.Println(name)
				return matrix, errors.New("image name doesn't exist for transitions")
			}

			startCol++
			columnName, err := excelize.ColumnNumberToName(startCol)
			if err != nil {
				return matrix, err
			}

			slowTimeSt, err := f.GetCellValue("Transitions", columnName+row)
			if err != nil {
				return matrix, err
			}

			slowTime, err := strconv.ParseFloat(slowTimeSt, 32)
			if err != nil {
				return matrix, err
			}

			startCol++
			columnName, err = excelize.ColumnNumberToName(startCol)
			if err != nil {
				return matrix, err
			}

			regTimeSt, err := f.GetCellValue("Transitions", columnName+row)
			if err != nil {
				return matrix, err
			}

			regTime, err := strconv.ParseFloat(regTimeSt, 32)
			if err != nil {
				return matrix, err
			}

			startCol++
			columnName, err = excelize.ColumnNumberToName(startCol)
			if err != nil {
				return matrix, err
			}

			fastTimeSt, err := f.GetCellValue("Transitions", columnName+row)
			if err != nil {
				return matrix, err
			}

			fastTime, err := strconv.ParseFloat(fastTimeSt, 32)
			if err != nil {
				return matrix, err
			}

			slowRep.ImageSetIDs = append(slowRep.ImageSetIDs, imageID)
			regRep.ImageSetIDs = append(regRep.ImageSetIDs, imageID)
			fastRep.ImageSetIDs = append(fastRep.ImageSetIDs, imageID)

			slowRep.Times = append(slowRep.Times, float32(slowTime))
			regRep.Times = append(regRep.Times, float32(regTime))
			fastRep.Times = append(fastRep.Times, float32(fastTime))

			slowRep.FullTime += float32(slowTime)
			regRep.FullTime += float32(regTime)
			fastRep.FullTime += float32(fastTime)

			startCol++
			columnName, err = excelize.ColumnNumberToName(startCol)
			if err != nil {
				return matrix, err
			}

			name, err = f.GetCellValue("Transitions", columnName+row)
			if err != nil {
				return matrix, err
			}
		}

		index1St, err := f.GetCellValue("Transitions", "A"+row)
		if err != nil {
			return matrix, err
		}

		index1, err := strconv.Atoi(index1St)
		if err != nil {
			return matrix, err
		}

		index2St, err := f.GetCellValue("Transitions", "B"+row)
		if err != nil {
			return matrix, err
		}

		index2, err := strconv.Atoi(index2St)
		if err != nil {
			return matrix, err
		}

		matrix.SlowMatrix[index1][index2] = slowRep
		matrix.RegularMatrix[index1][index2] = regRep
		matrix.FastMatrix[index1][index2] = fastRep

	}

	return matrix, nil
}
