package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setClusterCmd = &cobra.Command{
	Use:   "set-cluster",
	Short: "Set the current cluster",
	Long:  `Set the current cluster to connect to. This command updates the 'current-cluster' field in the configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run:   runSetCluster,
}

func init() {
	rootCmd.AddCommand(setClusterCmd)
}

func runSetCluster(cmd *cobra.Command, args []string) {
	clusterName := args[0]
	config := parseConfigFile()

	clusterExists := false
	for _, cluster := range config.Clusters {
		if cluster.Name == clusterName {
			clusterExists = true
			break
		}
	}

	if !clusterExists {
		fmt.Printf("Error: No cluster found with the name '%s' in the configuration.\n", clusterName)
		os.Exit(1)
	}

	viper.Set("current-cluster", clusterName)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	configFilePath := filepath.Join(homeDir, ".config", "esctl.yaml")
	err = viper.WriteConfigAs(configFilePath)
	if err != nil {
		fmt.Printf("Error writing updated configuration: %s\n", err)
		os.Exit(1)
	}
}
