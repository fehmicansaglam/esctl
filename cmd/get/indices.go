package get

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/cmd/config"
	"github.com/fehmicansaglam/esctl/cmd/utils"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var getIndicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "Get Elasticsearch indices",
	Long: utils.Trim(`
	Get Elasticsearch indices. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all indices.
	esctl get indices

	# Retrieve indices for a specific index.
	esctl get indices --index my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.ParseConfigFile()
		handleIndicesLogic(conf)
	},
}

func init() {
	getIndicesCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
}

var indexColumns = []output.ColumnDef{
	{Header: "INDEX", Type: output.Text},
	{Header: "UUID", Type: output.Text},
	{Header: "HEALTH", Type: output.Text},
	{Header: "STATUS", Type: output.Text},
	{Header: "SHARDS", Type: output.Number},
	{Header: "REPLICAS", Type: output.Number},
	{Header: "DOCS-COUNT", Type: output.Number},
	{Header: "DOCS-DELETED", Type: output.Number},
	{Header: "CREATION-DATE", Type: output.Date},
	{Header: "STORE-SIZE", Type: output.DataSize},
	{Header: "PRI-STORE-SIZE", Type: output.DataSize},
}

func handleIndicesLogic(conf config.Config) {
	indices, err := es.GetIndices(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve indices:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "index", indexColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, index := range indices {
		rowData := map[string]string{
			"INDEX":          index.Index,
			"UUID":           index.UUID,
			"HEALTH":         index.Health,
			"STATUS":         index.Status,
			"SHARDS":         index.Pri,
			"REPLICAS":       index.Rep,
			"DOCS-COUNT":     index.DocsCount,
			"DOCS-DELETED":   index.DocsDeleted,
			"CREATION-DATE":  index.CreationDate,
			"STORE-SIZE":     index.StoreSize,
			"PRI-STORE-SIZE": index.PriStoreSize,
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	if len(flagSortBy) > 0 {
		output.PrintTable(columnDefs, data, flagSortBy...)
	} else {
		output.PrintTable(columnDefs, data, "INDEX")
	}
}
