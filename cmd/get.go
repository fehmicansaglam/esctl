package cmd

import (
	"fmt"
	"os"
	"strconv"

	"text/tabwriter"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/shared"
	"github.com/spf13/cobra"
)

var (
	flagIndex        string
	flagShard        int
	flagPrimary      bool
	flagReplica      bool
	flagStarted      bool
	flagRelocating   bool
	flagInitializing bool
	flagUnassigned   bool
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
	shards, err := es.GetShards(shared.ElasticsearchHost, shared.ElasticsearchPort, flagIndex)

	if err != nil {
		panic(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "INDEX\tID\tSHARD\tPRI_REP\tSTATE\tDOCS\tSTORE\tIP\tNODE\tUNASSIGNED_REASON\tUNASSIGNED_AT\tSEGMENTS_COUNT")

	for _, shard := range shards {
		includeShardByState := false
		switch {
		case flagStarted && shard.State == constants.ShardStateStarted:
			includeShardByState = true
		case flagRelocating && shard.State == constants.ShardStateRelocating:
			includeShardByState = true
		case flagInitializing && shard.State == constants.ShardStateInitializing:
			includeShardByState = true
		case flagUnassigned && shard.State == constants.ShardStateUnassigned:
			includeShardByState = true
		case !flagStarted && !flagRelocating && !flagInitializing && !flagUnassigned:
			includeShardByState = true
		}

		shardNumber, err := strconv.Atoi(shard.Shard)
		if err != nil {
			panic(err)
		}
		includeShardByNumber := (flagShard == -1 || flagShard == shardNumber)

		includeShardByPriRep := (flagPrimary && shard.PriRep == constants.ShardPrimary) ||
			(flagReplica && shard.PriRep == constants.ShardReplica) ||
			(!flagPrimary && !flagReplica)

		if includeShardByState && includeShardByNumber && includeShardByPriRep {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				shard.Index, shard.ID, shard.Shard, humanizePriRep(shard.PriRep), shard.State, shard.Docs, shard.Store, shard.IP, shard.Node, shard.UnassignedReason, shard.UnassignedAt, shard.SegmentsCount)
		}
	}
}

func humanizePriRep(priRep string) string {
	switch priRep {
	case constants.ShardPrimary:
		return "primary"
	case constants.ShardReplica:
		return "replica"
	default:
		return priRep
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&flagIndex, "index", "", "Name of the index")
	getCmd.Flags().IntVar(&flagShard, "shard", -1, "Filter shards by shard number")
	getCmd.Flags().BoolVar(&flagPrimary, "primary", false, "Filter primary shards")
	getCmd.Flags().BoolVar(&flagReplica, "replica", false, "Filter replica shards")
	getCmd.Flags().BoolVar(&flagStarted, "started", false, "Filter shards in STARTED state")
	getCmd.Flags().BoolVar(&flagRelocating, "relocating", false, "Filter shards in RELOCATING state")
	getCmd.Flags().BoolVar(&flagInitializing, "initializing", false, "Filter shards in INITIALIZING state")
	getCmd.Flags().BoolVar(&flagUnassigned, "unassigned", false, "Filter shards in UNASSIGNED state")
}
