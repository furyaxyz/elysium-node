package nodebuilder

import (
	"github.com/elysiumorg/elysium-node/nodebuilder/das"
	"github.com/elysiumorg/elysium-node/nodebuilder/fraud"
	"github.com/elysiumorg/elysium-node/nodebuilder/header"
	"github.com/elysiumorg/elysium-node/nodebuilder/node"
	"github.com/elysiumorg/elysium-node/nodebuilder/p2p"
	"github.com/elysiumorg/elysium-node/nodebuilder/share"
	"github.com/elysiumorg/elysium-node/nodebuilder/state"
)

// PackageToAPI maps a package to its API struct. Currently only used for
// method discovery for openrpc spec generation
var PackageToAPI = map[string]interface{}{
	"fraud":  &fraud.API{},
	"state":  &state.API{},
	"share":  &share.API{},
	"header": &header.API{},
	"daser":  &das.API{},
	"p2p":    &p2p.API{},
	"node":   &node.API{},
}
