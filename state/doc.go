// Package state provides a structure for elysium-node's ability to access
// state-relevant information from as well as submit transactions/messages
// to the elysium network.
//
// This package contains one main interface, `Accessor`, that defines
// the methods available for both accessing and updating state on the
// elysium network.
//
// `Accessor` will contain three different implementations:
//  1. Implementation over a gRPC connection with a elysium-core node
//     called `CoreAccess`.
//  2. Implementation over a libp2p stream with a state-providing node.
//  3. Implementation over a local running instance of the
//     elysium-application (this feature will be implemented in *Full*
//     nodes).
package state
