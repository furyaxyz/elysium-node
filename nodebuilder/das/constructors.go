package das

import (
	"context"
	"fmt"

	"github.com/ipfs/go-datastore"

	"github.com/furyaxyz/go-fraud"
	libhead "github.com/furyaxyz/go-header"

	"github.com/furyaxyz/elysium-node/das"
	"github.com/furyaxyz/elysium-node/header"
	modfraud "github.com/furyaxyz/elysium-node/nodebuilder/fraud"
	"github.com/furyaxyz/elysium-node/share"
	"github.com/furyaxyz/elysium-node/share/eds/byzantine"
	"github.com/furyaxyz/elysium-node/share/p2p/shrexsub"
)

var _ Module = (*daserStub)(nil)

var errStub = fmt.Errorf("module/das: stubbed: dasing is not available on bridge nodes")

// daserStub is a stub implementation of the DASer that is used on bridge nodes, so that we can
// provide a friendlier error when users try to access the daser over the API.
type daserStub struct{}

func (d daserStub) SamplingStats(context.Context) (das.SamplingStats, error) {
	return das.SamplingStats{}, errStub
}

func (d daserStub) WaitCatchUp(context.Context) error {
	return errStub
}

func newDaserStub() Module {
	return &daserStub{}
}

func newDASer(
	da share.Availability,
	hsub libhead.Subscriber[*header.ExtendedHeader],
	store libhead.Store[*header.ExtendedHeader],
	batching datastore.Batching,
	fraudServ fraud.Service,
	bFn shrexsub.BroadcastFn,
	options ...das.Option,
) (*das.DASer, *modfraud.ServiceBreaker[*das.DASer], error) {
	ds, err := das.NewDASer(da, hsub, store, batching, fraudServ, bFn, options...)
	if err != nil {
		return nil, nil, err
	}

	return ds, &modfraud.ServiceBreaker[*das.DASer]{
		Service:   ds,
		FraudServ: fraudServ,
		FraudType: byzantine.BadEncoding,
	}, nil
}
