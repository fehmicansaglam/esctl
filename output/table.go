package output

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type ColumnDef struct {
	Header string
	Type   ColumnType
	Sort   func(i, j int) bool
}

type ColumnType int

const (
	Text ColumnType = iota
	Number
	Percent
	DataSize
	Date
)

func compareValues(left, right string, columnType ColumnType) bool {
	switch columnType {
	case Text:
		return sortText(left, right)
	case Number:
		return sortNumber(left, right)
	case DataSize:
		return sortDataSize(left, right)
	}
	return false
}

func PrintTable(columnDefs []ColumnDef, data [][]string, sortByHeaders ...string) {
	// Determine if a column is empty
	emptyColumns := make([]bool, len(columnDefs))
	for i := range columnDefs {
		empty := true
		for _, row := range data {
			if row[i] != "" {
				empty = false
				break
			}
		}
		emptyColumns[i] = empty
	}

	// Sort data if sortByHeaders are valid
	if len(sortByHeaders) > 0 {
		headerIndexMap := make(map[string]int)
		for i, columnDef := range columnDefs {
			headerIndexMap[strings.ToLower(columnDef.Header)] = i
		}

		sort.SliceStable(data, func(i, j int) bool {
			for _, header := range sortByHeaders {
				col, exists := headerIndexMap[strings.ToLower(header)]
				if !exists {
					continue
				}

				left, right := data[i][col], data[j][col]
				if left == right {
					continue
				}

				return compareValues(left, right, columnDefs[col].Type)
			}

			return false
		})
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Write headers
	for i, columnDef := range columnDefs {
		if !emptyColumns[i] {
			fmt.Fprintf(w, "%s\t", columnDef.Header)
		}
	}
	fmt.Fprintln(w)

	// Write data
	for _, row := range data {
		for i, cell := range row {
			if !emptyColumns[i] {
				fmt.Fprintf(w, "%s\t", cell)
			}
		}
		fmt.Fprintln(w)
	}
}
