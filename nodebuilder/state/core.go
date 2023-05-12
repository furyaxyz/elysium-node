package state

import (
	apptypes "github.com/elysiumorg/elysium-app/x/blob/types"
	libfraud "github.com/elysiumorg/go-fraud"
	"github.com/elysiumorg/go-header/sync"

	"github.com/elysiumorg/elysium-node/header"
	"github.com/elysiumorg/elysium-node/nodebuilder/core"
	modfraud "github.com/elysiumorg/elysium-node/nodebuilder/fraud"
	"github.com/elysiumorg/elysium-node/share/eds/byzantine"
	"github.com/elysiumorg/elysium-node/state"
)

// coreAccessor constructs a new instance of state.Module over
// a elysium-core connection.
func coreAccessor(
	corecfg core.Config,
	signer *apptypes.KeyringSigner,
	sync *sync.Syncer[*header.ExtendedHeader],
	fraudServ libfraud.Service,
) (*state.CoreAccessor, *modfraud.ServiceBreaker[*state.CoreAccessor]) {
	ca := state.NewCoreAccessor(signer, sync, corecfg.IP, corecfg.RPCPort, corecfg.GRPCPort)

	return ca, &modfraud.ServiceBreaker[*state.CoreAccessor]{
		Service:   ca,
		FraudType: byzantine.BadEncoding,
		FraudServ: fraudServ,
	}
}
