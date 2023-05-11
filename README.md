# esctl

esctl is a command-line interface (CLI) tool for easily retrieving read-only information about Elasticsearch clusters, including nodes, indices, and shards.

⚠️ **Warning: This tool is a work in progress and may not have full functionality or stability. Use it at your own risk.**

## Features

- Retrieve a list of all nodes in the Elasticsearch cluster.
- List all indices available in the Elasticsearch cluster.
- Get detailed information about shards, including their sizes and placement.

## Usage

esctl get nodes

Retrieves a list of all nodes in the Elasticsearch cluster.

esctl get indices

Retrieves a list of all indices in the Elasticsearch cluster.

esctl get shards

Retrieves detailed information about shards in the Elasticsearch cluster.

Please note that the 'get' command only provides read-only access and does not support data querying or modification operations.

## Installation

To install esctl, you need to have Go installed and set up. Then, you can use the following command:

```shell
go get github.com/fehmicansaglam/esctl
```

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
