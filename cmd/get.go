package cmd

import (
	"fmt"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/shared"
	"github.com/spf13/cobra"
)

var (
	indexName        string
	startedFlag      bool
	relocatingFlag   bool
	initializingFlag bool
	unassignedFlag   bool
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Elasticsearch entities",
	Long: `The 'get' command allows you to retrieve information about Elasticsearch entities. Supported entities include nodes, indices, and shards. This command provides a read-only view of the cluster and does not support data querying.

Usage:
	esctl get [entity]

Available Entities:
	- nodes: List all nodes in the Elasticsearch cluster.
	- indices: List all indices in the Elasticsearch cluster.
	- shards: List detailed information about shards, including their sizes and placement.

Options:
	[entity] - Specifies the entity type to retrieve. Supports 'nodes', 'indices', and 'shards'.

Examples:
	esctl get nodes
	Retrieves a list of all nodes in the Elasticsearch cluster.

	esctl get indices
	Retrieves a list of all indices in the Elasticsearch cluster.

	esctl get shards
	Retrieves detailed information about shards in the Elasticsearch cluster.

Please note that the 'get' command only provides read-only access and does not support data querying or modification operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify the entity (node, index, or shard).")
			return
		}

		entity := args[0]

		switch entity {
		case constants.EntityNode, constants.EntityNodes:
			// Retrieve and display information about Elasticsearch nodes
			// Your logic for handling the "node" entity goes here
			handleNodeLogic()
		case constants.EntityIndex, constants.EntityIndices:
			handleIndexLogic()
		case constants.EntityShard, constants.EntityShards:
			handleShardLogic()
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
			fmt.Println("Supported entities: node(s), index(es), shard(s)")

		}
	},
}

func handleNodeLogic() {
	// Logic for handling node-related functionality
	fmt.Println("Getting information about Elasticsearch nodes...")
}

func handleIndexLogic() {
	indices, err := es.GetIndices(shared.ElasticsearchHost, shared.ElasticsearchPort)
	if err != nil {
		panic(err)
	}
	for _, index := range indices {
		fmt.Println("Index\tHealth\tShards\tReplicas")
		fmt.Printf("%-9s(%s):\t%s\tshards,\t%s\treplicas. %s docs, %s\n", index.Index, index.Health, index.Shards, index.Replicas, index.DocsCount, index.StoreSize)
	}
}

func handleShardLogic() {
	shards, err := es.GetShards(shared.ElasticsearchHost, shared.ElasticsearchPort, indexName)

	if err != nil {
		panic(err)
	}

	for _, shard := range shards {
		includeShard := false

		switch {
		case startedFlag && shard.State == constants.ShardStateStarted:
			includeShard = true
		case relocatingFlag && shard.State == constants.ShardStateRelocating:
			includeShard = true
		case initializingFlag && shard.State == constants.ShardStateInitializing:
			includeShard = true
		case unassignedFlag && shard.State == constants.ShardStateUnassigned:
			includeShard = true
		case !startedFlag && !relocatingFlag && !initializingFlag && !unassignedFlag:
			includeShard = true
		}

		if includeShard {
			fmt.Println(shard.Index, shard.ID, shard.PriRep, shard.Shard, shard.IP, shard.State)
		}
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&indexName, "index", "", "Name of the index")
	getCmd.Flags().BoolVar(&startedFlag, "started", false, "Filter shards in STARTED state")
	getCmd.Flags().BoolVar(&relocatingFlag, "relocating", false, "Filter shards in RELOCATING state")
	getCmd.Flags().BoolVar(&initializingFlag, "initializing", false, "Filter shards in INITIALIZING state")
	getCmd.Flags().BoolVar(&unassignedFlag, "unassigned", false, "Filter shards in UNASSIGNED state")
}
