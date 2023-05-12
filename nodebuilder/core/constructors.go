package core

import (
	"github.com/elysiumorg/elysium-node/core"
)

func remote(cfg Config) (core.Client, error) {
	return core.NewRemote(cfg.IP, cfg.RPCPort)
}
