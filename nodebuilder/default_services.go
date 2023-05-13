package nodebuilder

import (
	"github.com/furyaxyz/elysium-node/nodebuilder/das"
	"github.com/furyaxyz/elysium-node/nodebuilder/fraud"
	"github.com/furyaxyz/elysium-node/nodebuilder/header"
	"github.com/furyaxyz/elysium-node/nodebuilder/node"
	"github.com/furyaxyz/elysium-node/nodebuilder/p2p"
	"github.com/furyaxyz/elysium-node/nodebuilder/share"
	"github.com/furyaxyz/elysium-node/nodebuilder/state"
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
