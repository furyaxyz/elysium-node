package gateway

import (
	"github.com/furyaxyz/elysium-node/api/gateway"
	"github.com/furyaxyz/elysium-node/das"
	"github.com/furyaxyz/elysium-node/nodebuilder/header"
	"github.com/furyaxyz/elysium-node/nodebuilder/share"
	"github.com/furyaxyz/elysium-node/nodebuilder/state"
)

// Handler constructs a new RPC Handler from the given services.
func Handler(
	state state.Module,
	share share.Module,
	header header.Module,
	daser *das.DASer,
	serv *gateway.Server,
) {
	handler := gateway.NewHandler(state, share, header, daser)
	handler.RegisterEndpoints(serv)
	handler.RegisterMiddleware(serv)
}

func server(cfg *Config) *gateway.Server {
	return gateway.NewServer(cfg.Address, cfg.Port)
}
