package share

import (
	disc "github.com/furyaxyz/elysium-node/share/availability/discovery"
	"github.com/furyaxyz/elysium-node/share/getters"
	"github.com/furyaxyz/elysium-node/share/p2p/peers"
	"github.com/furyaxyz/elysium-node/share/p2p/shrexeds"
	"github.com/furyaxyz/elysium-node/share/p2p/shrexnd"
)

// WithPeerManagerMetrics is a utility function to turn on peer manager metrics and that is
// expected to be "invoked" by the fx lifecycle.
func WithPeerManagerMetrics(m *peers.Manager) error {
	return m.WithMetrics()
}

// WithDiscoveryMetrics is a utility function to turn on discovery metrics and that is expected to
// be "invoked" by the fx lifecycle.
func WithDiscoveryMetrics(d *disc.Discovery) error {
	return d.WithMetrics()
}

func WithShrexClientMetrics(edsClient *shrexeds.Client, ndClient *shrexnd.Client) error {
	err := edsClient.WithMetrics()
	if err != nil {
		return err
	}

	return ndClient.WithMetrics()
}

func WithShrexServerMetrics(edsServer *shrexeds.Server, ndServer *shrexnd.Server) error {
	err := edsServer.WithMetrics()
	if err != nil {
		return err
	}

	return ndServer.WithMetrics()
}

func WithShrexGetterMetrics(sg *getters.ShrexGetter) error {
	return sg.WithMetrics()
}
