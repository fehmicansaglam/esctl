package cmd

import (
	"fmt"
	"os"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/es"
	"github.com/fehmicansaglam/esctl/output"
	"github.com/fehmicansaglam/esctl/shared"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe <entity>",
	Short: "Print detailed information about an entity",
	Long:  "Print detailed information about the specified entity.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]
		switch entity {
		case constants.EntityCluster:
			handleDescribeCluster()
		// Handle other entities here
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
			os.Exit(1)
		}
	},
}

func handleDescribeCluster() {
	cluster, err := es.GetCluster(shared.ElasticsearchHost, shared.ElasticsearchPort)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintYaml(cluster)
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
