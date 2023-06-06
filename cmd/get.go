package cmd

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndex        string
	flagNode         string
	flagShard        int
	flagPrimary      bool
	flagReplica      bool
	flagStarted      bool
	flagRelocating   bool
	flagInitializing bool
	flagUnassigned   bool
	flagActions      []string
	flagSortBy       []string
)

var getCmd = &cobra.Command{
	Use:   "get ENTITY",
	Short: "Get Elasticsearch entities",
	Long: trim(`
The 'get' command allows you to retrieve information about Elasticsearch entities.

Available Entities:
  - nodes: List all nodes in the Elasticsearch cluster.
  - indices: List all indices in the Elasticsearch cluster.
  - shards: List detailed information about shards, including their sizes and placement.
  - aliases: List all aliases in the Elasticsearch cluster.
  - tasks: List all tasks in the Elasticsearch cluster.`),
	Example: trimAndIndent(`
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
	Args:       cobra.ExactArgs(1),
	ValidArgs:  []string{"node", "index", "shard", "alias", "task"},
	ArgAliases: []string{"nodes", "indices", "shards", "aliases", "tasks"},
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]

		config := parseConfigFile()

		switch entity {
		case constants.EntityNode, constants.EntityNodes:
			handleNodeLogic(config)
		case constants.EntityIndex, constants.EntityIndices:
			handleIndexLogic(config)
		case constants.EntityShard, constants.EntityShards:
			handleShardLogic(config)
		case constants.EntityAlias, constants.EntityAliases:
			handleAliasLogic(config)
		case constants.EntityTask, constants.EntityTasks:
			handleTaskLogic(config)
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
			os.Exit(1)
		}
	},
}

func getColumnDefs(config Config, entity string, defaultColumns []output.ColumnDef) []output.ColumnDef {
	entityConfig, ok := config.Entities[entity]
	if !ok || len(entityConfig.Columns) == 0 {
		return defaultColumns
	}

	columnDefs := make([]output.ColumnDef, 0, len(entityConfig.Columns))
	for _, configColumn := range entityConfig.Columns {
		var found bool
		for _, defaultColumn := range defaultColumns {
			if defaultColumn.Header == configColumn {
				columnDefs = append(columnDefs, defaultColumn)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "Unknown column: %s\n", configColumn)
			os.Exit(1)
		}
	}

	return columnDefs
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&flagIndex, "index", "", "Name of the index")
	getCmd.Flags().StringVar(&flagNode, "node", "", "Filter shards by node")
	getCmd.Flags().IntVar(&flagShard, "shard", -1, "Filter shards by shard number")
	getCmd.Flags().BoolVar(&flagPrimary, "primary", false, "Filter primary shards")
	getCmd.Flags().BoolVar(&flagReplica, "replica", false, "Filter replica shards")
	getCmd.Flags().BoolVar(&flagStarted, "started", false, "Filter shards in STARTED state")
	getCmd.Flags().BoolVar(&flagRelocating, "relocating", false, "Filter shards in RELOCATING state")
	getCmd.Flags().BoolVar(&flagInitializing, "initializing", false, "Filter shards in INITIALIZING state")
	getCmd.Flags().BoolVar(&flagUnassigned, "unassigned", false, "Filter shards in UNASSIGNED state")
	getCmd.Flags().StringSliceVar(&flagActions, "actions", []string{}, "Filter tasks by actions")
	getCmd.Flags().StringSliceVar(&flagSortBy, "sort-by", []string{}, "Columns to sort by (comma-separated)")
}
