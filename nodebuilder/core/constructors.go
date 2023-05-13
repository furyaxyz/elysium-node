package core

import (
	"github.com/furyaxyz/elysium-node/core"
)

func remote(cfg Config) (core.Client, error) {
	return core.NewRemote(cfg.IP, cfg.RPCPort)
}
