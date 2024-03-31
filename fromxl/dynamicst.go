package fromxl

import (
	"errors"
	"fmt"
	"i9posesa/assets"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDynamics(imageSetMap map[string]string) ([]assets.DynamicStr, error) {
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

	dynamics := []assets.DynamicStr{}

	for i := 3; i < 1000; i++ {
		name, err := f.GetCellValue("DynamicStretches", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if name == "" {
			break
		}

		secsStr, err := f.GetCellValue("DynamicStretches", "C"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if secsStr == "" {
			posList2, err := getPosListDynamic(strconv.Itoa(i), f, imageSetMap)
			if err != nil {
				return nil, err
			}

			dynamics[len(dynamics)-1].PositionSlice2 = posList2

		} else {
			secs, err := strconv.ParseFloat(secsStr, 32)
			if err != nil {
				return nil, err
			}

			separates, err := f.GetCellValue("DynamicStretches", "D"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			backendID, err := f.GetCellValue("DynamicStretches", "E"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			posList1, err := getPosListDynamic(strconv.Itoa(i), f, imageSetMap)
			if err != nil {
				return nil, err
			}

			dynamics = append(dynamics, assets.DynamicStr{
				Name:           name,
				Secs:           float32(secs),
				SeparateSets:   separates != "",
				BackendID:      backendID,
				PositionSlice1: posList1,
			})

			id, err := f.GetCellValue("DynamicStretches", "A"+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}

			if id != "" {

				primID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return nil, err
				}
				dynamics[len(dynamics)-1].ID = primID
			}

		}

	}

	return dynamics, nil
}

func getPosListDynamic(row string, f *excelize.File, imageSetMap map[string]string) ([]assets.StrPosition, error) {

	ret := []assets.StrPosition{}

	startCol := 6

	columnName, err := excelize.ColumnNumberToName(startCol)
	if err != nil {
		return nil, err
	}

	hasMore, err := f.GetCellValue("DynamicStretches", columnName+row)
	if err != nil {
		return nil, err
	}

	for hasMore != "" {

		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		posName, err := f.GetCellValue("DynamicStretches", columnName+row)
		if err != nil {
			return nil, err
		}

		imageSetID, ok := imageSetMap[posName]
		if !ok {
			fmt.Println(columnName + row)
			return nil, errors.New("image set name not in existing list in db")
		}

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		percentSecsStr, err := f.GetCellValue("DynamicStretches", columnName+row)
		if err != nil {
			return nil, err
		}

		percentSecs, err := strconv.ParseFloat(percentSecsStr, 32)
		if err != nil {
			return nil, err
		}

		ret = append(ret, assets.StrPosition{
			ImageSetID:  imageSetID,
			PercentSecs: float32(percentSecs),
		})

		startCol++
		columnName, err = excelize.ColumnNumberToName(startCol)
		if err != nil {
			return nil, err
		}

		hasMore, err = f.GetCellValue("DynamicStretches", columnName+row)
		if err != nil {
			return nil, err
		}

	}

	return ret, nil
}
