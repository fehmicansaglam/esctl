package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagTerm    []string
	flagExists  []string
	flagGroupBy string
)

var countCmd = &cobra.Command{
	Use:   "count [INDEX]",
	Short: "Count documents in an index or in all indices",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var index string
		if len(args) == 1 {
			index = args[0]
		}
		handleCount(index)
	},
}

func handleCount(index string) {
	var counts map[string]es.GroupCount
	var err error

	counts, err = es.CountDocuments(flagIndex, flagTerm, flagExists, flagGroupBy)
	if err != nil {
		fmt.Printf("Failed to get document counts: %v\n", err)
		os.Exit(1)
	}

	columnDefs := []output.ColumnDef{
		{Header: "INDEX", Type: output.Text},
		{Header: strings.ToUpper(flagGroupBy), Type: output.Text},
		{Header: "COUNT", Type: output.Number},
	}

	data := [][]string{}

	for index, groupCount := range counts {
		for group, count := range groupCount {
			rowData := map[string]string{
				"INDEX":                      index,
				strings.ToUpper(flagGroupBy): group,
				"COUNT":                      strconv.Itoa(count),
			}

			row := make([]string, len(columnDefs))
			for i, colDef := range columnDefs {
				row[i] = rowData[colDef.Header]
			}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(columnDefs, data, flagSortBy...)
	} else {
		output.PrintTable(columnDefs, data, "INDEX")
	}
}

func init() {
	countCmd.Flags().StringSliceVar(&flagTerm, "term", []string{}, "Term filters to apply")
	countCmd.Flags().StringSliceVar(&flagExists, "exists", []string{}, "Exists filters to apply")
	countCmd.Flags().StringVar(&flagGroupBy, "group-by", "", "Field to group the documents by")

	rootCmd.AddCommand(countCmd)
}
