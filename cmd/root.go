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
	setupElasticsearchProtocol()
	setupElasticsearchUsername()
	setupElasticsearchPassword()
	setupElasticsearchHost()
	setupElasticsearchPort()
}

func setupElasticsearchProtocol() {
	defaultProtocol := constants.DefaultElasticsearchProtocol
	defaultProtocolEnv := os.Getenv(constants.ElasticsearchProtocolEnvVar)
	if defaultProtocolEnv != "" {
		defaultProtocol = defaultProtocolEnv
	}
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchProtocol, "protocol", defaultProtocol, "Elasticsearch protocol")
}

func setupElasticsearchUsername() {
	defaultUsername := os.Getenv(constants.ElasticsearchUsernameEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchUsername, "username", defaultUsername, "Elasticsearch username")
}

func setupElasticsearchPassword() {
	defaultPassword := os.Getenv(constants.ElasticsearchPasswordEnvVar)
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchPassword, "password", defaultPassword, "Elasticsearch password")
}

func setupElasticsearchHost() {
	defaultHost := os.Getenv(constants.ElasticsearchHostEnvVar)
	if defaultHost == "" {
		defaultHost = constants.DefaultElasticsearchHost
	}
	rootCmd.PersistentFlags().StringVar(&shared.ElasticsearchHost, "host", defaultHost, "Elasticsearch host")
}

func setupElasticsearchPort() {
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
