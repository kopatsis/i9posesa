package fromxl

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func GetImageSetNameMap() (map[string]string, error) {
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

	ret := map[string]string{}

	for i := 3; i < 1200; i++ {
		name, err := f.GetCellValue("ImageSets", "B"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		if name == "" {
			break
		}

		id, err := f.GetCellValue("ImageSets", "C"+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		ret[name] = id

	}

	return ret, nil
}
