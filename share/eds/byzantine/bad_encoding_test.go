package byzantine

import (
	"context"
	"testing"

	mdutils "github.com/ipfs/go-merkledag/test"
	"github.com/stretchr/testify/require"
	core "github.com/tendermint/tendermint/types"

	"github.com/elysiumorg/elysium-app/pkg/da"
	"github.com/elysiumorg/rsmt2d"

	"github.com/elysiumorg/elysium-node/header"
	"github.com/elysiumorg/elysium-node/share"
	"github.com/elysiumorg/elysium-node/share/ipld"
)

// TestIncorrectBadEncodingFraudProof asserts that BEFP is not generated for the correct data
func TestIncorrectBadEncodingFraudProof(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bServ := mdutils.Bserv()

	squareSize := 8
	shares := share.RandShares(t, squareSize*squareSize)

	eds, err := share.AddShares(ctx, shares, bServ)
	require.NoError(t, err)

	dah := da.NewDataAvailabilityHeader(eds)

	// get an arbitrary row
	row := uint(squareSize / 2)
	rowShares := eds.Row(row)
	rowRoot := dah.RowsRoots[row]

	shareProofs, err := GetProofsForShares(ctx, bServ, ipld.MustCidFromNamespacedSha256(rowRoot), rowShares)
	require.NoError(t, err)

	// create a fake error for data that was encoded correctly
	fakeError := ErrByzantine{
		Index:  uint32(row),
		Shares: shareProofs,
		Axis:   rsmt2d.Row,
	}

	h := &header.ExtendedHeader{
		RawHeader: core.Header{
			Height: 420,
		},
		DAH: &dah,
		Commit: &core.Commit{
			BlockID: core.BlockID{
				Hash: []byte("made up hash"),
			},
		},
	}

	proof := CreateBadEncodingProof(h.Hash(), uint64(h.Height()), &fakeError)
	err = proof.Validate(h)
	require.Error(t, err)
}
