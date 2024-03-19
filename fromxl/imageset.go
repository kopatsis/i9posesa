package fromxl

import (
	"fmt"
	"i9posesa/assets"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var arraySize = 4

func GetImageSets() ([]assets.ImageSet, error) {
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

	imageSets := []assets.ImageSet{}

	for i := 3; i < 1000; i++ {
		name, err := f.GetCellValue("ImageSets", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if name == "" {
			break
		}

		current := assets.ImageSet{
			Name: name,
		}

		id, err := f.GetCellValue("ImageSets", "A"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if id != "" {
			primID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, err
			}
			current.ID = primID
		}

		start, low := 3, []string{}
		for i := 0; i < arraySize; i++ {

			columnName, err := excelize.ColumnNumberToName(start + i)
			if err != nil {
				return nil, err
			}
			val, err := f.GetCellValue("ImageSets", columnName+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}
			low = append(low, val)

		}
		current.Low = low

		start, mid := 3+arraySize, []string{}
		for i := 0; i < arraySize; i++ {

			columnName, err := excelize.ColumnNumberToName(start + i)
			if err != nil {
				return nil, err
			}
			val, err := f.GetCellValue("ImageSets", columnName+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}
			mid = append(mid, val)

		}
		current.Mid = mid

		start, high := 3+(2*arraySize), []string{}
		for i := 0; i < arraySize; i++ {

			columnName, err := excelize.ColumnNumberToName(start + i)
			if err != nil {
				return nil, err
			}
			val, err := f.GetCellValue("ImageSets", columnName+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}
			high = append(high, val)

		}
		current.High = high

		start, original := 3+(3*arraySize), []string{}
		for i := 0; i < arraySize; i++ {

			columnName, err := excelize.ColumnNumberToName(start + i)
			if err != nil {
				return nil, err
			}
			val, err := f.GetCellValue("ImageSets", columnName+strconv.Itoa(i))
			if err != nil {
				return nil, err
			}
			original = append(original, val)

		}
		current.Original = original

		imageSets = append(imageSets, current)

	}

	return imageSets, nil
}
