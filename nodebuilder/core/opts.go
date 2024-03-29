package core

import (
	"go.uber.org/fx"

	"github.com/furyaxyz/elysium-node/core"
	"github.com/furyaxyz/elysium-node/header"
	"github.com/furyaxyz/elysium-node/libs/fxutil"
)

// WithClient sets custom client for core process
func WithClient(client core.Client) fx.Option {
	return fxutil.ReplaceAs(client, new(core.Client))
}

// WithHeaderConstructFn sets custom func that creates extended header
func WithHeaderConstructFn(construct header.ConstructFn) fx.Option {
	return fx.Replace(construct)
}
