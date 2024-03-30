package toxl

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WriteToXL(sheet string, data []primitive.ObjectID) error {
	f, err := excelize.OpenFile("assets/posesaxl.xlsx")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		if err := f.Save(); err != nil {
			fmt.Println(err)
		}
	}()

	currentRow := 2
	for _, id := range data {
		hasID, err := f.GetCellValue(sheet, "A"+strconv.Itoa(currentRow))
		if err != nil {
			return err
		}

		for hasID != "" {
			currentRow++
			hasID, err = f.GetCellValue(sheet, "A"+strconv.Itoa(currentRow))
			if err != nil {
				return err
			}
		}

		err = f.SetCellStr(sheet, "A"+strconv.Itoa(currentRow), id.Hex())
		if err != nil {
			return err
		}
		currentRow++

	}

	return nil
}

func WriteTransitionXL(data []primitive.ObjectID) error {

	if len(data) == 0 {
		return nil
	}

	f, err := excelize.OpenFile("assets/posesaxl.xlsx")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		if err := f.Save(); err != nil {
			fmt.Println(err)
		}
	}()

	err = f.SetCellStr("Transitions", "B1", data[0].Hex())
	if err != nil {
		return err
	}

	return nil
}

func WriteToXLDblRow(sheet string, data []primitive.ObjectID) error {
	f, err := excelize.OpenFile("assets/posesaxl.xlsx")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() {
		if err := f.Save(); err != nil {
			fmt.Println(err)
		}
	}()

	currentRow := 2
	for _, id := range data {
		hasID, err := f.GetCellValue(sheet, "A"+strconv.Itoa(currentRow))
		if err != nil {
			return err
		}

		hasThirdRow, err := f.GetCellValue(sheet, "C"+strconv.Itoa(currentRow))
		if err != nil {
			return err
		}

		for hasID != "" || hasThirdRow == "" {
			currentRow++
			hasID, err = f.GetCellValue(sheet, "A"+strconv.Itoa(currentRow))
			if err != nil {
				return err
			}
			hasThirdRow, err = f.GetCellValue(sheet, "C"+strconv.Itoa(currentRow))
			if err != nil {
				return err
			}
		}

		err = f.SetCellStr(sheet, "A"+strconv.Itoa(currentRow), id.Hex())
		if err != nil {
			return err
		}
		currentRow++

	}

	return nil
}
