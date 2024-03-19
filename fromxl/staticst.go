package fromxl

import (
	"errors"
	"fmt"
	"i9posesa/assets"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetStatics(imageSetMap map[string]string) ([]assets.StaticStr, error) {
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

	statics := []assets.StaticStr{}

	for i := 3; i < 1000; i++ {
		name, err := f.GetCellValue("StaticStretches", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if name == "" {
			break
		}

		posName, err := f.GetCellValue("StaticStretches", "C"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		imageSetID, ok := imageSetMap[posName]
		if !ok {
			return nil, errors.New("image set name not in existing list in db")
		}

		statics = append(statics, assets.StaticStr{
			Name:        name,
			ImageSetID1: imageSetID,
		})

		id, err := f.GetCellValue("StaticStretches", "A"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if id != "" {

			primID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, err
			}
			statics[len(statics)-1].ID = primID
		}

		posName2, err := f.GetCellValue("StaticStretches", "D"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if posName2 != "" {
			imageSetID, ok = imageSetMap[posName]
			if !ok {
				return nil, errors.New("image set name not in existing list in db")
			}
			statics[len(statics)-1].ImageSetID2 = imageSetID
		}

	}

	return statics, nil
}
