package output

import (
	"fmt"
	"strconv"
	"strings"
)

func sortText(left, right string) bool {
	return left < right
}

func sortNumber(left, right string) bool {
	a, _ := strconv.ParseFloat(left, 64)
	b, _ := strconv.ParseFloat(right, 64)
	return a < b
}

func sortDataSize(left, right string) bool {
	size1, _ := parseDataSize(left)
	size2, _ := parseDataSize(right)
	return size1 < size2
}

func parseDataSize(sizeStr string) (float64, error) {
	sizeStr = strings.ToLower(sizeStr)
	var value float64
	var unit string
	var err error

	for i := 0; i < len(sizeStr); i++ {
		if (sizeStr[i] < '0' || sizeStr[i] > '9') && sizeStr[i] != '.' {
			value, err = strconv.ParseFloat(sizeStr[:i], 64)
			if err != nil {
				fmt.Println(err)
				return 0, err
			}
			unit = sizeStr[i:]
			break
		}
	}

	switch unit {
	case "b":
		return value, nil
	case "kb":
		return value * 1024, nil
	case "mb":
		return value * 1024 * 1024, nil
	case "gb":
		return value * 1024 * 1024 * 1024, nil
	case "tb":
		return value * 1024 * 1024 * 1024 * 1024, nil
	default:
		fmt.Printf("unknown unit: %s\n", unit)
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}
}
