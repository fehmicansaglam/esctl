package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fehmicansaglam/esctl/constants"
	"github.com/fehmicansaglam/esctl/shared"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl is CLI for Elasticsearch",
	Long:  `esctl is a read-only CLI for Elasticsearch that allows users to manage and monitor their Elasticsearch clusters.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initialize)

	initProtocol()
	initHost()
	initPort()
	initUsername()
	initPassword()
}

func initialize() {
	if shared.ElasticsearchHost == "" {
		config := parseConfigFile()
		readClusterFromConfig(config)
	}
}

func readClusterFromConfig(config Config) {
	if len(config.Clusters) == 0 {
		fmt.Println("Error: No clusters defined in the configuration.")
		os.Exit(1)
	}
	if config.CurrentCluster == "" {
		config.CurrentCluster = config.Clusters[0].Name
	}

	clusterFound := false
	for _, cluster := range config.Clusters {
		if cluster.Name == config.CurrentCluster {
			shared.ElasticsearchProtocol = cluster.Protocol
			if shared.ElasticsearchProtocol == "" {
				shared.ElasticsearchProtocol = constants.DefaultElasticsearchProtocol
			}
			shared.ElasticsearchPort = cluster.Port
			if shared.ElasticsearchPort == 0 {
				shared.ElasticsearchPort = constants.DefaultElasticsearchPort
			}
			shared.ElasticsearchUsername = cluster.Username
			shared.ElasticsearchPassword = cluster.Password
			shared.ElasticsearchHost = cluster.Host
			if shared.ElasticsearchHost == "" {
				fmt.Println("Error: 'host' field is not specified in the configuration for the current cluster.")
				os.Exit(1)
			}
			clusterFound = true
			break
		}
	}

	if !clusterFound {
		fmt.Printf("Error: No cluster found with the name '%s' in the configuration.\n", config.CurrentCluster)
		os.Exit(1)
	}
}

func initProtocol() {
	defaultProtocol := constants.DefaultElasticsearchProtocol
	defaultProtocolEnv := os.Getenv(constants.ElasticsearchProtocolEnvVar)
	if defaultProtocolEnv != "" {
		defaultProtocol = defaultProtocolEnv
	}
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchProtocol, "protocol", defaultProtocol, "Elasticsearch protocol")
}

func initHost() {
	defaultHost := os.Getenv(constants.ElasticsearchHostEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchHost, "host", defaultHost, "Elasticsearch host")
}

func initPort() {
	defaultPort := constants.DefaultElasticsearchPort
	defaultPortStr := os.Getenv(constants.ElasticsearchPortEnvVar)
	if defaultPortStr != "" {
		parsedPort, err := strconv.Atoi(defaultPortStr)
		if err != nil || parsedPort <= 0 {
			fmt.Printf("Invalid value for %s environment variable: %s\n", constants.ElasticsearchPortEnvVar, defaultPortStr)
			os.Exit(1)
		}
		defaultPort = parsedPort
	}
	rootCmd.PersistentFlags().IntVar(&shared.ElasticsearchPort, "port", defaultPort, "Elasticsearch port")
}

func initUsername() {
	defaultUsername := os.Getenv(constants.ElasticsearchUsernameEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchUsername, "username", defaultUsername, "Elasticsearch username")
}

func initPassword() {
	defaultPassword := os.Getenv(constants.ElasticsearchPasswordEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchPassword, "password", defaultPassword, "Elasticsearch password")
}
