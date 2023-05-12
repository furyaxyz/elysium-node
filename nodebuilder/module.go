package nodebuilder

import (
	"context"

	"go.uber.org/fx"

	"github.com/elysiumorg/elysium-node/libs/fxutil"
	"github.com/elysiumorg/elysium-node/nodebuilder/core"
	"github.com/elysiumorg/elysium-node/nodebuilder/das"
	"github.com/elysiumorg/elysium-node/nodebuilder/fraud"
	"github.com/elysiumorg/elysium-node/nodebuilder/gateway"
	"github.com/elysiumorg/elysium-node/nodebuilder/header"
	"github.com/elysiumorg/elysium-node/nodebuilder/node"
	"github.com/elysiumorg/elysium-node/nodebuilder/p2p"
	"github.com/elysiumorg/elysium-node/nodebuilder/rpc"
	"github.com/elysiumorg/elysium-node/nodebuilder/share"
	"github.com/elysiumorg/elysium-node/nodebuilder/state"
)

func ConstructModule(tp node.Type, network p2p.Network, cfg *Config, store Store) fx.Option {
	log.Infow("Accessing keyring...")
	ks, err := store.Keystore()
	if err != nil {
		fx.Error(err)
	}
	signer, err := state.KeyringSigner(cfg.State, ks, network)
	if err != nil {
		fx.Error(err)
	}

	baseComponents := fx.Options(
		fx.Supply(tp),
		fx.Supply(network),
		fx.Provide(p2p.BootstrappersFor),
		fx.Provide(func(lc fx.Lifecycle) context.Context {
			return fxutil.WithLifecycle(context.Background(), lc)
		}),
		fx.Supply(cfg),
		fx.Supply(store.Config),
		fx.Provide(store.Datastore),
		fx.Provide(store.Keystore),
		fx.Supply(node.StorePath(store.Path())),
		fx.Supply(signer),
		// modules provided by the node
		p2p.ConstructModule(tp, &cfg.P2P),
		state.ConstructModule(tp, &cfg.State),
		header.ConstructModule(tp, &cfg.Header),
		share.ConstructModule(tp, &cfg.Share),
		rpc.ConstructModule(tp, &cfg.RPC),
		gateway.ConstructModule(tp, &cfg.Gateway),
		core.ConstructModule(tp, &cfg.Core),
		das.ConstructModule(tp, &cfg.DASer),
		fraud.ConstructModule(tp),
		node.ConstructModule(tp),
	)

	return fx.Module(
		"node",
		baseComponents,
	)
}
