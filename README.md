# esctl

esctl is a command-line interface (CLI) tool for easily retrieving read-only information about Elasticsearch clusters, including nodes, indices, and shards.

⚠️ **Warning: This tool is a work in progress and may not have full functionality or stability. Use it at your own risk.**

## Features

- Retrieve a list of all nodes in the Elasticsearch cluster.
- List all indices available in the Elasticsearch cluster.
- Get detailed information about shards, including their sizes and placement.

## Contributing

Thank you for your interest in contributing to esctl! While I don't accept pull requests at the moment, I encourage you to open issues for bug reports, feature requests, or any other suggestions you may have. Your feedback helps me improve the tool.

When opening an issue, please provide as much information as possible, including steps to reproduce the problem or a detailed description of the requested feature. This will help me better understand and address the issue or request.

I aim to have monthly releases where I address the reported issues and include requested features. While I can't guarantee immediate fixes, your reported issues will be considered for upcoming releases. I appreciate your patience and understanding in this regard.

I value your contributions and feedback, and I'm grateful for your support in making esctl better!

## Usage

esctl get nodes

Retrieves a list of all nodes in the Elasticsearch cluster.

esctl get indices

Retrieves a list of all indices in the Elasticsearch cluster.

### Get Shards

To retrieve shards from Elasticsearch, you can use the following command:

```shell
esctl get shards [--index <index_name>] [--started] [--relocating] [--initializing] [--unassigned]
```

* `--index <index_name>`: Specifies the name of the index to retrieve shards from.
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

Please note that the 'get' command only provides read-only access and does not support data querying or modification operations.

## Installation

To install `esctl`, ensure that you have Go installed and set up in your development environment. Then, follow the steps below:

1. Open a terminal or command prompt.

2. Run the following command to install `esctl`:

   ```shell
   go install github.com/fehmicansaglam/esctl
   ```
   This command will fetch the source code from the GitHub repository, compile it, and install the `esctl` binary in your Go workspace's `bin` directory.

3. Make sure that your Go workspace's `bin` directory is added to your system's `PATH` environment variable. This step will allow you to run `esctl` from any directory in the terminal or command prompt.

Once installed, you can run `esctl` by simply typing `esctl` in the terminal or command prompt.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

