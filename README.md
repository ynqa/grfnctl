# grfnctl

*grfnctl* is an unofficial CLI for
[Grafana API](https://grafana.com/docs/grafana/latest/developers/http_api/).
It helps you manage Grafana resources such as dashboards, datasources,
and snapshots efficiently from the terminal.

## Features

- Dashboard
    - [x] Create or update
    - [x] List
    - [x] Export as JSON
- Datasource
    - [x] List
    - [x] (Test to query)
- Snapshot
    - [x] Create ([Ceveats below](#caveats-create-snapshot))
    - [x] Delete
    - [x] List
- Role
    - [x] Display current user information
- Auto-completion for commands and resource identifiers

## Installation

### Homebrew

```bash
brew install ynqa/tap/grfnctl
```

### Go

```bash
go install github.com/ynqa/grfnctl/cmd/grfnctl@latest
```

## Prerequisite

Export your Grafana URL and API token as environment variables:

```bash
# your Grafana instance URL
export GRAFANA_URL="https://your-grafana-instance.com"

# token with appropriate permissions
export GRAFANA_TOKEN="your_api_token"
# or basic auth
export GRAFANA_USER="your_username"
export GRAFANA_PASSWORD="your_password"
```

## Usage

```bash
Grafana CLI tool

Usage:
  grfnctl [command]

Available Commands:
  completion  Generate shell completion script
  dashboard   Provide Grafana dashboard related commands
  datasource  Provide Grafana data source related commands
  help        Help about any command
  snapshot    Provide Grafana snapshots related commands
  whoami      Display the current user information

Flags:
  -h, --help   help for grfnctl

Use "grfnctl [command] --help" for more information about a command.

## Version

Display the CLI version with:

```bash
grfnctl --version
```

When building from source, override the embedded version string via:

```bash
go build -ldflags "-X github.com/ynqa/grfnctl/cmd.version=vX.Y.Z" ./cmd/grfnctl
```

## Auto-completion

You can enable shell auto-completion for *grfnctl* commands. For example, to enable bash completion, run:

```bash
source <(grfnctl completion bash)
# or zsh
source <(grfnctl completion zsh)
```

And also, you can add the above command to your shell profile
(e.g., `~/.bashrc` or `~/.zshrc`) to enable auto-completion automatically on shell startup.

## Caveats: Create Snapshot

Grafana exposes an API for creating snapshots
([ref](https://grafana.com/docs/grafana/v12.2/developers/http_api/snapshot/#create-new-snapshot)).

That endpoint assumes clients submit the complete payloads of query results verbatim
(in addition, some fields must be converted) as snapshot JSON.
*grfnctl* converts a provided dashboard JSON into a snapshot JSON
and sends it to the API, but because no official conversion recipe exists,
the implementation relies on observed behavior rather than documented guarantees.

As of v0.1.0, *grfnctl* supports only the following:
- Data source
  - [x] Prometheus
- Data format
  - [x] timeseries

Snapshot creation has not been validated for every scenario, so it may fail under certain conditions.
