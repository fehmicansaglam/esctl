package cmd

import (
	"os"

	"github.com/fehmicansaglam/esctl/shared"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl is CLI for Elasticsearch",
	Long: `esctl is a read-only Command Line Interface tool for Elasticsearch that allows
users to manage and monitor their Elasticsearch clusters.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchHost, "host", "localhost", "Elasticsearch host")
	rootCmd.PersistentFlags().IntVar(&shared.ElasticsearchPort, "port", 9200, "Elasticsearch port")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
