package tabular

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func PrintTable(headers []string, data [][]string) {
	// Determine if a column is empty
	emptyColumns := make([]bool, len(headers))
	for i := range headers {
		empty := true
		for _, row := range data {
			if row[i] != "" {
				empty = false
				break
			}
		}
		emptyColumns[i] = empty
	}

	// Create a tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Write headers
	for i, h := range headers {
		if !emptyColumns[i] {
			fmt.Fprintf(w, "%s\t", h)
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
