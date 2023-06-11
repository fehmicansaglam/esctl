package cmd

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
)

var aliasColumns = []output.ColumnDef{
	{Header: "ALIAS", Type: output.Text},
	{Header: "INDEX", Type: output.Text},
}

func handleAliasLogic(config Config) {
	aliases, err := es.GetAliases(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve aliases:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(config, "alias", aliasColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for alias, index := range aliases {
		rowData := map[string]string{
			"ALIAS": alias,
			"INDEX": index,
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	output.PrintTable(columnDefs, data, flagSortBy...)
}
