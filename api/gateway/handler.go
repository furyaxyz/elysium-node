package gateway

import (
	logging "github.com/ipfs/go-log/v2"

	"github.com/furyaxyz/elysium-node/das"
	"github.com/furyaxyz/elysium-node/nodebuilder/header"
	"github.com/furyaxyz/elysium-node/nodebuilder/share"
	"github.com/furyaxyz/elysium-node/nodebuilder/state"
)

var log = logging.Logger("gateway")

type Handler struct {
	state  state.Module
	share  share.Module
	header header.Module
	das    *das.DASer
}

func NewHandler(
	state state.Module,
	share share.Module,
	header header.Module,
	das *das.DASer,
) *Handler {
	return &Handler{
		state:  state,
		share:  share,
		header: header,
		das:    das,
	}
}
