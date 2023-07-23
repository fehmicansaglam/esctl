package get

import (
	"fmt"
	"strings"

	"github.com/fehmicansaglam/esctl/cmd/config"
	"github.com/fehmicansaglam/esctl/cmd/utils"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Elasticsearch entities",
	Long: utils.Trim(`
The 'get' command allows you to retrieve information about Elasticsearch entities.

Available Entities:
  - nodes: List all nodes in the Elasticsearch cluster.
  - indices: List all indices in the Elasticsearch cluster.
  - shards: List detailed information about shards, including their sizes and placement.
  - aliases: List all aliases in the Elasticsearch cluster.
  - tasks: List all tasks in the Elasticsearch cluster.`),
	Example: utils.TrimAndIndent(`
#Retrieve a list of all nodes in the Elasticsearch cluster.
esctl get nodes

#Retrieve a list of all indices in the Elasticsearch cluster.
esctl get indices

#Retrieve detailed information about shards in the Elasticsearch cluster.
esctl get shards

#Retrieve shard information for an index.
esctl get shards --index my_index

#Retrieve shard information filtered by state.
esctl get shards --started --relocating

#Retrieve all aliases.
esctl get aliases

#Retrieve tasks filtered by actions using wildcard patterns.
esctl get tasks --actions 'index*' --actions '*search*'

#Retrieve all tasks.
esctl get tasks`),
}

func init() {
	getCmd.PersistentFlags().StringSliceVarP(&flagSortBy, "sort-by", "s", []string{}, "Columns to sort by (comma-separated)")
	getCmd.PersistentFlags().StringSliceVarP(&flagColumns, "columns", "c", []string{}, "Columns to display (comma-separated) or 'all'")

	getCmd.AddCommand(getAliasesCmd)
	getCmd.AddCommand(getIndicesCmd)
	getCmd.AddCommand(getNodesCmd)
	getCmd.AddCommand(getShardsCmd)
	getCmd.AddCommand(getTasksCmd)
}

func Cmd() *cobra.Command {
	return getCmd
}

func buildColumnDefs(columns []string, defaultColumns []output.ColumnDef) ([]output.ColumnDef, error) {
	columnDefs := make([]output.ColumnDef, 0, len(columns))
	for _, column := range columns {
		var found bool
		for _, defaultColumn := range defaultColumns {
			if strings.EqualFold(defaultColumn.Header, column) {
				columnDefs = append(columnDefs, defaultColumn)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unknown column: %s", column)
		}
	}
	return columnDefs, nil
}

func getColumnDefs(conf config.Config, entity string, defaultColumns []output.ColumnDef) ([]output.ColumnDef, error) {
	if len(flagColumns) > 0 {
		for _, column := range flagColumns {
			if strings.EqualFold(column, "all") {
				return defaultColumns, nil
			}
		}
		return buildColumnDefs(flagColumns, defaultColumns)
	} else {
		entityConfig, ok := conf.Entities[entity]
		if !ok || len(entityConfig.Columns) == 0 {
			return defaultColumns, nil
		}
		return buildColumnDefs(entityConfig.Columns, defaultColumns)
	}
}
