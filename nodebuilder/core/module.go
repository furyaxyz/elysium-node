package core

import (
	"context"

	"go.uber.org/fx"

	libhead "github.com/elysiumorg/go-header"

	"github.com/elysiumorg/elysium-node/core"
	"github.com/elysiumorg/elysium-node/header"
	"github.com/elysiumorg/elysium-node/libs/fxutil"
	"github.com/elysiumorg/elysium-node/nodebuilder/node"
	"github.com/elysiumorg/elysium-node/nodebuilder/p2p"
	"github.com/elysiumorg/elysium-node/share/eds"
	"github.com/elysiumorg/elysium-node/share/p2p/shrexsub"
)

// ConstructModule collects all the components and services related to managing the relationship
// with the Core node.
func ConstructModule(tp node.Type, cfg *Config, options ...fx.Option) fx.Option {
	// sanitize config values before constructing module
	cfgErr := cfg.Validate()

	baseComponents := fx.Options(
		fx.Supply(*cfg),
		fx.Error(cfgErr),
		fx.Options(options...),
	)

	switch tp {
	case node.Light, node.Full:
		return fx.Module("core", baseComponents)
	case node.Bridge:
		return fx.Module("core",
			baseComponents,
			fx.Provide(core.NewBlockFetcher),
			fxutil.ProvideAs(core.NewExchange, new(libhead.Exchange[*header.ExtendedHeader])),
			fx.Invoke(fx.Annotate(
				func(
					bcast libhead.Broadcaster[*header.ExtendedHeader],
					fetcher *core.BlockFetcher,
					pubsub *shrexsub.PubSub,
					construct header.ConstructFn,
					store *eds.Store,
				) *core.Listener {
					return core.NewListener(bcast, fetcher, pubsub.Broadcast, construct, store, p2p.BlockTime)
				},
				fx.OnStart(func(ctx context.Context, listener *core.Listener) error {
					return listener.Start(ctx)
				}),
				fx.OnStop(func(ctx context.Context, listener *core.Listener) error {
					return listener.Stop(ctx)
				}),
			)),
			fx.Provide(fx.Annotate(
				remote,
				fx.OnStart(func(ctx context.Context, client core.Client) error {
					return client.Start()
				}),
				fx.OnStop(func(ctx context.Context, client core.Client) error {
					return client.Stop()
				}),
			)),
		)
	default:
		panic("invalid node type")
	}
}
