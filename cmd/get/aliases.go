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

var getAliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Get Elasticsearch aliases",
	Long: utils.Trim(`
	Get Elasticsearch aliases. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all aliases.
	esctl get aliases

	# Retrieve aliases for a specific index.
	esctl get aliases --index my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.ParseConfigFile()
		handleAliasLogic(conf)
	},
}

func init() {
	getAliasesCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
}

var aliasColumns = []output.ColumnDef{
	{Header: "ALIAS", Type: output.Text},
	{Header: "INDEX", Type: output.Text},
}

func handleAliasLogic(conf config.Config) {
	aliases, err := es.GetAliases(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve aliases:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "alias", aliasColumns)
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
