# Elysium Node

[![Go Reference](https://pkg.go.dev/badge/github.com/furyaxyz/elysium-node.svg)](https://pkg.go.dev/github.com/furyaxyz/elysium-node)
[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/furyaxyz/elysium-node)](https://github.com/furyaxyz/elysium-node/releases/latest)
[![Go CI](https://github.com/furyaxyz/elysium-node/actions/workflows/go-ci.yml/badge.svg)](https://github.com/furyaxyz/elysium-node/actions/workflows/go-ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/furyaxyz/elysium-node)](https://goreportcard.com/report/github.com/furyaxyz/elysium-node)
[![codecov](https://codecov.io/gh/furyaxyz/elysium-node/branch/main/graph/badge.svg?token=CWGA4RLDS9)](https://codecov.io/gh/furyaxyz/elysium-node)

Golang implementation of Elysium's data availability node types (`light` | `full` | `bridge`).

The elysium-node types described above comprise the elysium data availability (DA) network.

The DA network wraps the elysium-core consensus network by listening for blocks from the consensus network and making them digestible for data availability sampling (DAS).

Continue reading [here](https://blog.elysium.org/elysium-mvp-release-data-availability-sampling-light-clients) if you want to learn more about DAS and how it enables secure and scalable access to Elysium chain data.

## Table of Contents

- [Elysium Node](#elysium-node)
  - [Table of Contents](#table-of-contents)
  - [Minimum requirements](#minimum-requirements)
  - [System Requirements](#system-requirements)
  - [Installation](#installation)
  - [API docs](#api-docs)
  - [Node types](#node-types)
  - [Run a node](#run-a-node)
    - [Environment variables](#environment-variables)
  - [Package-specific documentation](#package-specific-documentation)
  - [Code of Conduct](#code-of-conduct)

## Minimum requirements

| Requirement | Notes          |
| ----------- |----------------|
| Go version  | 1.20 or higher |

## System Requirements

See the official docs page for system requirements per node type:

- [Bridge](https://docs.elysium.org/nodes/bridge-node#hardware-requirements)
- [Light](https://docs.elysium.org/nodes/light-node#hardware-requirements)
- [Full](https://docs.elysium.org/nodes/full-storage-node#hardware-requirements)

## Installation

```sh
git clone https://github.com/furyaxyz/elysium-node.git
cd elysium-node
make build
sudo make install
```

For more information on setting up a node and the hardware requirements needed, go visit our docs at <https://docs.elysium.org>.

## API docs

Elysium-node public API is documented [here](https://docs.elysium.org/category/node-api/).

## Node types

- **Bridge** nodes - relay blocks from the elysium consensus network to the elysium data availability (DA) network
- **Full** nodes - fully reconstruct and store blocks by sampling the DA network for shares
- **Light** nodes - verify the availability of block data by sampling the DA network for shares

More information can be found [here](https://github.com/furyaxyz/elysium-node/blob/main/docs/adr/adr-003-march2022-testnet.md#legend).

## Run a node

`<node_type>` can be `bridge`, `full` or `light`.

```sh
elysium <node_type> init
```

```sh
elysium <node_type> start
```

### Environment variables

| Variable                | Explanation                         | Default value | Required |
| ----------------------- | ----------------------------------- | ------------- | -------- |
| `ELYSIUM_BOOTSTRAPPER` | Start the node in bootstrapper mode | `false`       | Optional |

## Package-specific documentation

- [Header](./header/doc.go)
- [Share](./share/doc.go)
- [DAS](./das/doc.go)

## Code of Conduct

See our Code of Conduct [here](https://docs.elysium.org/community/coc).
