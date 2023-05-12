package getters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

	"github.com/elysiumorg/nmt"
	"github.com/elysiumorg/nmt/namespace"

	"github.com/elysiumorg/elysium-node/libs/utils"
	"github.com/elysiumorg/elysium-node/share"
	"github.com/elysiumorg/elysium-node/share/ipld"
)

var (
	tracer = otel.Tracer("share/getters")
	log    = logging.Logger("share/getters")

	errOperationNotSupported = errors.New("operation is not supported")
)

// filterRootsByNamespace returns the row roots from the given share.Root that contain the passed
// namespace ID.
func filterRootsByNamespace(root *share.Root, nID namespace.ID) []cid.Cid {
	rowRootCIDs := make([]cid.Cid, 0, len(root.RowsRoots))
	for _, row := range root.RowsRoots {
		if !nID.Less(nmt.MinNamespace(row, nID.Size())) && nID.LessOrEqual(nmt.MaxNamespace(row, nID.Size())) {
			rowRootCIDs = append(rowRootCIDs, ipld.MustCidFromNamespacedSha256(row))
		}
	}
	return rowRootCIDs
}

// collectSharesByNamespace collects NamespaceShares within the given namespace ID from the given
// share.Root.
func collectSharesByNamespace(
	ctx context.Context,
	bg blockservice.BlockGetter,
	root *share.Root,
	nID namespace.ID,
) (shares share.NamespacedShares, err error) {
	ctx, span := tracer.Start(ctx, "collect-shares-by-namespace", trace.WithAttributes(
		attribute.String("root", root.String()),
		attribute.String("nid", nID.String()),
	))
	defer func() {
		utils.SetStatusAndEnd(span, err)
	}()

	rootCIDs := filterRootsByNamespace(root, nID)
	if len(rootCIDs) == 0 {
		return nil, share.ErrNotFound
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	shares = make([]share.NamespacedRow, len(rootCIDs))
	for i, rootCID := range rootCIDs {
		// shadow loop variables, to ensure correct values are captured
		i, rootCID := i, rootCID
		errGroup.Go(func() error {
			row, proof, err := share.GetSharesByNamespace(ctx, bg, rootCID, nID, len(root.RowsRoots))
			shares[i] = share.NamespacedRow{
				Shares: row,
				Proof:  proof,
			}
			if err != nil {
				return fmt.Errorf("retrieving nID %x for row %x: %w", nID, rootCID, err)
			}
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	// return ErrNotFound if no shares are found for namespaceID
	if len(rootCIDs) == 1 && len(shares[0].Shares) == 0 {
		return nil, share.ErrNotFound
	}

	return shares, nil
}

func verifyNIDSize(nID namespace.ID) error {
	if len(nID) != share.NamespaceSize {
		return fmt.Errorf("expected namespace ID of size %d, got %d",
			share.NamespaceSize, len(nID))
	}
	return nil
}

// ctxWithSplitTimeout will split timeout stored in context by splitFactor and return the result if
// it is greater than minTimeout. minTimeout == 0 will be ignored, splitFactor <= 0 will be ignored
func ctxWithSplitTimeout(
	ctx context.Context,
	splitFactor int,
	minTimeout time.Duration,
) (context.Context, context.CancelFunc) {
	deadline, ok := ctx.Deadline()
	if !ok || splitFactor <= 0 {
		if minTimeout == 0 {
			return context.WithCancel(ctx)
		}
		return context.WithTimeout(ctx, minTimeout)
	}

	timeout := time.Until(deadline) / time.Duration(splitFactor)
	if minTimeout == 0 || timeout > minTimeout {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithTimeout(ctx, minTimeout)
}

// ErrorContains reports whether any error in err's tree matches any error in targets tree.
func ErrorContains(err, target error) bool {
	if errors.Is(err, target) || target == nil {
		return true
	}

	target = errors.Unwrap(target)
	if target == nil {
		return false
	}
	return ErrorContains(err, target)
}
