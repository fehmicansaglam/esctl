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

var listClustersCmd = &cobra.Command{
	Use:   "list-clusters",
	Short: "List the clusters defined in the esctl.yml file",
	Run:   runListClusters,
}

func init() {
	rootCmd.AddCommand(setClusterCmd)
	rootCmd.AddCommand(listClustersCmd)
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

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing updated configuration: %s\n", err)
		os.Exit(1)
	}
}

func runListClusters(cmd *cobra.Command, args []string) {
	config := parseConfigFile()
	for _, cluster := range config.Clusters {
		clusterName := cluster.Name
		if clusterName == config.CurrentCluster {
			clusterName += "(*)"
		}
		fmt.Printf("- name: %s\n", clusterName)
		fmt.Printf("  host: %s\n", cluster.Host)
		if cluster.Protocol != "" {
			fmt.Printf("  protocol: %s\n", cluster.Protocol)
		}
		if cluster.Port != 0 {
			fmt.Printf("  port: %d\n", cluster.Port)
		}
		if cluster.Username != "" {
			fmt.Printf("  username: %s\n", cluster.Username)
		}
		if cluster.Password != "" {
			fmt.Printf("  password: %s\n", cluster.Password)
		}
	}
}

type Cluster struct {
	Name     string `mapstructure:"name"`
	Protocol string `mapstructure:"protocol"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Entity struct {
	Columns []string `mapstructure:"columns"`
}

type Config struct {
	CurrentCluster string            `mapstructure:"current-cluster"`
	Clusters       []Cluster         `mapstructure:"clusters"`
	Entities       map[string]Entity `mapstructure:"entities"`
}

func parseConfigFile() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user's home directory: %v\n", err)
		os.Exit(1)
	}

	viper.AddConfigPath(filepath.Join(home, ".config"))
	viper.SetConfigName("esctl")
	viper.SetConfigType("yml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Error unmarshaling config into struct: %v\n", err)
		os.Exit(1)
	}

	return config
}
