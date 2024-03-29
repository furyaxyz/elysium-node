package state

import (
	apptypes "github.com/furyaxyz/elysium-app/x/blob/types"
	libfraud "github.com/furyaxyz/go-fraud"
	"github.com/furyaxyz/go-header/sync"

	"github.com/furyaxyz/elysium-node/header"
	"github.com/furyaxyz/elysium-node/nodebuilder/core"
	modfraud "github.com/furyaxyz/elysium-node/nodebuilder/fraud"
	"github.com/furyaxyz/elysium-node/share/eds/byzantine"
	"github.com/furyaxyz/elysium-node/state"
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
