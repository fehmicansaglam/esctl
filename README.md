# esctl

[![Latest Stable Version](https://img.shields.io/badge/version-v0.1.0-blue.svg)](https://github.com/fehmicansaglam/esctl/releases/tag/v0.1.0)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/fehmicansaglam/esctl/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/fehmicansaglam/esctl)](https://goreportcard.com/report/github.com/fehmicansaglam/esctl)

`esctl` is a command-line tool for managing and interacting with Elasticsearch clusters. It provides a convenient interface to perform various operations, such as querying cluster information, managing indices, retrieving shard details, and monitoring tasks.

## Features:
- Retrieve information about nodes, indices, shards, aliases, and tasks in an Elasticsearch cluster
- Filter and sort data based on various criteria
- Perform read-only operations for cluster monitoring and analysis
- Simple and intuitive command-line interface

## Contributing
Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## Table of Contents

- [Installation](#installation)
- [Examples](#examples)
- [Usage](#usage)
  - [Elasticsearch Host Configuration](#elasticsearch-host-configuration)
  - [Get](#get)
  - [Describe](#describe)
- [License](#license)

## Installation

To install `esctl`, ensure that you have Go installed and set up in your development environment. Then, follow the steps below:

1. Open a terminal or command prompt.

2. Run the following command to install `esctl`:

   ```shell
   go install github.com/fehmicansaglam/esctl@latest
   ```
   This command will fetch the source code from the GitHub repository, compile it, and install the `esctl` binary in your Go workspace's `bin` directory.

3. Make sure that your Go workspace's `bin` directory is added to your system's `PATH` environment variable. This step will allow you to run `esctl` from any directory in the terminal or command prompt.

Once installed, you can run `esctl` by simply typing `esctl` in the terminal or command prompt.

## Examples

```shell
> esctl get shards --index=articles --primary
INDEX     ID                      SHARD  PRI-REP  STATE    DOCS  STORE  IP         NODE               SEGMENTS-COUNT
articles  jxn-Oa3XSPigaCBYt9fKiw  0      primary  STARTED  0     225b   127.0.0.1  es-data-0          0
articles  jxn-Oa3XSPigaCBYt9fKiw  1      primary  STARTED  0     225b   127.0.0.1  es-data-0          0
articles  jxn-Oa3XSPigaCBYt9fKiw  2      primary  STARTED  0     225b   127.0.0.1  es-data-0          0
```

```shell
> esctl get shards --index=articles --shard 0 --unassigned
INDEX     SHARD  PRI-REP  STATE       UNASSIGNED-REASON  UNASSIGNED-AT
articles  0      replica  UNASSIGNED  CLUSTER_RECOVERED  2023-05-07T20:37:07.520Z
articles  0      replica  UNASSIGNED  CLUSTER_RECOVERED  2023-05-07T20:37:07.520Z
```

```shell
> esctl get indices
HEALTH  STATUS  INDEX     UUID                    PRI  REP  DOCS-COUNT  DOCS-DELETED  CREATION-DATE             STORE-SIZE  PRI-STORE-SIZE
yellow  open    articles  8vCars4rQquYHNhpKV2fow  3    2    0           0             2023-05-07T19:17:52.259Z  675b        675b
```

```shell
> esctl get nodes
NAME               IP         NODE-ROLE    MASTER  HEAP-MAX  HEAP-CURRENT  HEAP-PERCENT  CPU  LOAD-1M  DISK-TOTAL  DISK-USED  DISK-AVAILABLE
es-data-0          127.0.0.1  cdfhilmrstw  *       4gb       1.6gb         41%           10%  2.02     232.9gb     199.2gb    33.6gb
```

```shell
> esctl get aliases --index=articles
ALIAS           INDEX
articles_alias  articles
```
## Usage

### Elasticsearch Host Configuration

`esctl` allows you to configure the Elasticsearch host and port using the `--host` and `--port` flags or the `ELASTICSEARCH_HOST` and `ELASTICSEARCH_PORT` environment variables. By default, the host is set to `localhost` and the port is set to `9200`.

To specify a custom host, you can use the `--host` flag followed by the desired host value. For example:

```shell
esctl --host=<your_host> <command>
```

Similarly, to specify a custom port, you can use the `--port` flag followed by the desired port value. For example:

```shell
esctl --port=<your_port> <command>
```

Alternatively, you can set the `ELASTICSEARCH_HOST` and `ELASTICSEARCH_PORT` environment variables to your desired Elasticsearch host and port, respectively. If the `--host` and `--port` flags are not provided and the corresponding environment variables are set, `esctl` will use the values from the environment variables as the host and port.

If the `--host` and `--port` flags are not provided and the `ELASTICSEARCH_HOST` and `ELASTICSEARCH_PORT` environment variables are not set, `esctl` will default to `localhost` and `9200`, respectively.

### Get

The `get` command allows you to retrieve information about Elasticsearch entities. Supported entities include nodes, indices, shards, aliases, and tasks. This command provides a read-only view of the cluster and does not support data querying.

```shell
esctl get [entity] [flags]
```

#### Available Entities

- `nodes`: List all nodes in the Elasticsearch cluster.
- `indices`: List all indices in the Elasticsearch cluster.
- `shards`: List detailed information about shards, including their sizes and placement.
- `aliases`: List all aliases in the Elasticsearch cluster.
- `tasks`: List all tasks in the Elasticsearch cluster.

#### Flags

- `--index`: Specifies the name of the index (applies to `indices` and `shards` entities).
- `--shard`: Filters shards by shard number (applies to `shards` entity).
- `--primary`: Filters primary shards (applies to `shards` entity).
- `--replica`: Filters replica shards (applies to `shards` entity).
- `--started`: Filters shards in STARTED state (applies to `shards` entity).
- `--relocating`: Filters shards in RELOCATING state (applies to `shards` entity).
- `--initializing`: Filters shards in INITIALIZING state (applies to `shards` entity).
- `--unassigned`: Filters shards in UNASSIGNED state (applies to `shards` entity).
- `--actions`: Filters tasks by actions (applies to `tasks` entity).
- `--sort-by`: Specifies the columns to sort by, separated by commas (applies to all entities).


#### Get Nodes

Retrieves a list of all nodes in the Elasticsearch cluster.

```shell
esctl get nodes
```

#### Get Indices

Retrieves a list of all indices in the Elasticsearch cluster.

```shell
esctl get indices
```

#### Get Shards

To retrieve shards from Elasticsearch, you can use the following command:

```shell
esctl get shards [--index <index_name>] [--shard <shard>] [--primary] [--replica] [--started] [--relocating] [--initializing] [--unassigned]
```

* `--index <index_name>`: Specifies the name of the index to retrieve shards from.
* `--shard <shard>`: Filters shards by shard number.
* `--primary`: Filters primary shards.
* `--replica`: Filters replica shards.
* `--started`: Filters shards in the STARTED state.
* `--relocating`: Filters shards in the RELOCATING state.
* `--initializing`: Filters shards in the INITIALIZING state.
* `--unassigned`: Filters shards in the UNASSIGNED state.

If none of the flags are provided, all shards will be returned.

Example usage:

```shell
esctl get shards --index my_index --relocating
```
This will retrieve only the shards that are currently relocating for the specified index.

#### Get Aliases

Retrieves the list of aliases defined in Elasticsearch, including the index names they are associated with.

Usage:

```shell
esctl get aliases [--index <index_name>]
```

Options:

`--index`: (optional) Filter the aliases by a specific index. If not provided, aliases from all indices will be returned.

#### Get Tasks

The `get tasks` command retrieves information about tasks in the Elasticsearch cluster.

Usage:

```shell
esctl get tasks [--actions <actions>...]
```

Example:

```shell
esctl get tasks --actions 'index*' --actions '*search*'

```
### Describe

The `esctl describe` command allows you to retrieve detailed information about various entities in the Elasticsearch cluster. The output is in YAML format, making it easy to read and understand.

#### Describe Cluster

This command outputs the cluster information in YAML format, providing a comprehensive overview of the cluster's current state.

```shell
esctl describe cluster
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

