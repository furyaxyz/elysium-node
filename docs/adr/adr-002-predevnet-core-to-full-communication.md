# ADR #002: Devnet Elysium Core <> Elysium Node Communication

## Authors

@renaynay @Wondertan

## Changelog

* 2021-09-09: initial draft

## Legend

**Elysium Core** = tendermint consensus node that lives in [the elysium-core repository](https://github.com/furyaxyz/elysium-core).

**Elysium Node** = elysium `full` or `light` nodes that live in [this repository](https://github.com/furyaxyz/elysium-node).

## Context

After the offsite, there was a bit of confusion on what the default behaviour for running a Elysium Full node should be in the devnet. Since we decided on an architecture where Core nodes will communicate with Elysium Full nodes *exclusively* via RPC (over an HTTP connection, for example), it is necessary that a certain percentage of Elysium Full nodes in the devnet run with either an embedded Core node process or are able to fetch block information from a remote Core endpoint, otherwise there would be no way for the two separate networks (Core network and Elysium network) to communicate.

## Decision

Since the flow of information in devnet is unidirectional, where Core nodes provide block information to Elysium Full nodes, the default behaviour for running a Elysium Full node is to have an embedded Core node process running within the Full node itself. Not only will this ensure that at least some Elysium Full nodes in the network will be communicating with Core nodes, it also makes it easier for end users to spin up a Elysium Full node without having to worry about feeding the Elysium Full node a remote Core endpoint from which it would fetch information.

It is also important to note that for devnet, it should also be possible to run Elysium Full nodes as `standalone` processes (without a trusted remote or embedded Core node) as Elysium Full nodes should also be capable of learning of block information on a P2P-level from other Elysium Full nodes.

## Detailed Design

* Elysium Full node should be able to run in a `default` mode where Elysium Full node embeds a Core node process
* Elysium Full node should also be able to be started with a `--core.disable` flag to indicate that the Elysium Full node will be running *without* a Core node process
* Elysium Full node should be able to take in a `--core.remote` endpoint that would indicate to the Full node that it should *not* embed the Core node process, but rather dial the provided remote Core node endpoint.
* Elysium Full nodes that rely on Core node processes (whether embedded or remote) should also communicate with other Elysium Full nodes on a P2P-level, broadcasting new headers from blocks that they've fetched from the Core nodes *and* being able to handle broadcasted block-related messages from other Full nodes on the network.

It is preferable that a devnet-ready Elysium Full node is *agnostic* to the method by which it receives new block information. Therefore, we will abstract the interface related to "fetching blocks" so that in the view of the Elysium Full node, it does not care *how* it is receiving blocks, only that it *is* receiving new blocks.

## Consequences of embedding Elysium Core process into Elysium Full node

### Positive

* Better UX for average devnet users who do not want to deal with spinning up a Elysium Core node and passing the endpoint to the Elysium Full node.
* Makes it easier to guarantee that there will be *some* Full nodes in the devnet that will be fetching blocks from Elysium Core nodes.

### Negative

* Eventually this work will be rendered useless as communicating with Elysium Core over RPC is a crutch we decided to use in order to streamline interoperability between Core and Full nodes. All communication beyond devnet will be over the P2P layer.

## Status

Proposed
