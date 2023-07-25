package output

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var chunkifyRegexp = regexp.MustCompile(`(\d+|\D+)`)

// Chunkify splits a string into chunks of numeric and non-numeric substrings using a regular expression.
func chunkify(s string) []string {
	return chunkifyRegexp.FindAllString(s, -1)
}

// Returns true if the first string precedes the second one according to natural order.
func sortText(left, right string) bool {
	chunksLeft := chunkify(left)
	chunksRight := chunkify(right)

	for i := range chunksLeft {
		if i >= len(chunksRight) {
			// If we reached the last chunk of right, then left is greater than right.
			return false
		}

		if chunksLeft[i] != chunksRight[i] {
			leftInt, leftErr := strconv.Atoi(chunksLeft[i])
			rightInt, rightErr := strconv.Atoi(chunksRight[i])

			if leftErr == nil && rightErr == nil {
				// If both chunks are numeric, compare them as integers
				return leftInt < rightInt
			}

			// If one or both chunks are non-numeric, perform lexicographic comparison
			return chunksLeft[i] < chunksRight[i]
		}
	}

	// All chunks are equal
	return true
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

func sortPercent(val1, val2 string) bool {
	parsedVal1, _ := strconv.ParseFloat(strings.TrimRight(val1, "%"), 64)
	parsedVal2, _ := strconv.ParseFloat(strings.TrimRight(val2, "%"), 64)
	return parsedVal1 > parsedVal2
}

func sortDate(val1, val2 string) bool {
	time1, err := time.Parse("2006-01-02T15:04:05.999Z", val1)
	if err != nil {
		return false
	}

	time2, err := time.Parse("2006-01-02T15:04:05.999Z", val2)
	if err != nil {
		return false
	}

	return time1.Before(time2)
}

func parseDataSize(sizeStr string) (float64, error) {
	if sizeStr == "" {
		return 0, nil
	}

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
