package getters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"

	"github.com/furyaxyz/nmt/namespace"
	"github.com/furyaxyz/rsmt2d"

	"github.com/furyaxyz/elysium-node/share"
	"github.com/furyaxyz/elysium-node/share/p2p"
	"github.com/furyaxyz/elysium-node/share/p2p/peers"
	"github.com/furyaxyz/elysium-node/share/p2p/shrexeds"
	"github.com/furyaxyz/elysium-node/share/p2p/shrexnd"
)

var _ share.Getter = (*ShrexGetter)(nil)

const (
	// defaultMinRequestTimeout value is set according to observed time taken by healthy peer to
	// serve getEDS request for block size 256
	defaultMinRequestTimeout = time.Minute // should be >= shrexeds server write timeout
	defaultMinAttemptsCount  = 3
)

var meter = global.MeterProvider().Meter("shrex/getter")

type metrics struct {
	edsAttempts syncint64.Histogram
	ndAttempts  syncint64.Histogram
}

func (m *metrics) recordEDSAttempt(ctx context.Context, attemptCount int, success bool) {
	if m == nil {
		return
	}
	if ctx.Err() != nil {
		ctx = context.Background()
	}
	m.edsAttempts.Record(ctx, int64(attemptCount), attribute.Bool("success", success))
}

func (m *metrics) recordNDAttempt(ctx context.Context, attemptCount int, success bool) {
	if m == nil {
		return
	}
	if ctx.Err() != nil {
		ctx = context.Background()
	}
	m.ndAttempts.Record(ctx, int64(attemptCount), attribute.Bool("success", success))
}

func (sg *ShrexGetter) WithMetrics() error {
	edsAttemptHistogram, err := meter.SyncInt64().Histogram(
		"getters_shrex_eds_attempts_per_request",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of attempts per shrex/eds request"),
	)
	if err != nil {
		return err
	}

	ndAttemptHistogram, err := meter.SyncInt64().Histogram(
		"getters_shrex_nd_attempts_per_request",
		instrument.WithUnit(unit.Dimensionless),
		instrument.WithDescription("Number of attempts per shrex/nd request"),
	)
	if err != nil {
		return err
	}

	sg.metrics = &metrics{
		edsAttempts: edsAttemptHistogram,
		ndAttempts:  ndAttemptHistogram,
	}
	return nil
}

// ShrexGetter is a share.Getter that uses the shrex/eds and shrex/nd protocol to retrieve shares.
type ShrexGetter struct {
	edsClient *shrexeds.Client
	ndClient  *shrexnd.Client

	peerManager *peers.Manager

	// minRequestTimeout limits minimal timeout given to single peer by getter for serving the request.
	minRequestTimeout time.Duration
	// minAttemptsCount will be used to split request timeout into multiple attempts. It will allow to
	// attempt multiple peers in scope of one request before context timeout is reached
	minAttemptsCount int

	metrics *metrics
}

func NewShrexGetter(edsClient *shrexeds.Client, ndClient *shrexnd.Client, peerManager *peers.Manager) *ShrexGetter {
	return &ShrexGetter{
		edsClient:         edsClient,
		ndClient:          ndClient,
		peerManager:       peerManager,
		minRequestTimeout: defaultMinRequestTimeout,
		minAttemptsCount:  defaultMinAttemptsCount,
	}
}

func (sg *ShrexGetter) Start(ctx context.Context) error {
	return sg.peerManager.Start(ctx)
}

func (sg *ShrexGetter) Stop(ctx context.Context) error {
	return sg.peerManager.Stop(ctx)
}

func (sg *ShrexGetter) GetShare(context.Context, *share.Root, int, int) (share.Share, error) {
	return nil, fmt.Errorf("getter/shrex: GetShare %w", errOperationNotSupported)
}

func (sg *ShrexGetter) GetEDS(ctx context.Context, root *share.Root) (*rsmt2d.ExtendedDataSquare, error) {
	var (
		attempt int
		err     error
	)
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		attempt++
		start := time.Now()
		peer, setStatus, getErr := sg.peerManager.Peer(ctx, root.Hash())
		if getErr != nil {
			err = errors.Join(err, getErr)
			log.Debugw("eds: couldn't find peer",
				"hash", root.String(),
				"err", getErr,
				"finished (s)", time.Since(start))
			sg.metrics.recordEDSAttempt(ctx, attempt, false)
			return nil, fmt.Errorf("getter/shrex: %w", err)
		}

		reqStart := time.Now()
		reqCtx, cancel := ctxWithSplitTimeout(ctx, sg.minAttemptsCount-attempt+1, sg.minRequestTimeout)
		eds, getErr := sg.edsClient.RequestEDS(reqCtx, root.Hash(), peer)
		cancel()
		switch {
		case getErr == nil:
			setStatus(peers.ResultSynced)
			sg.metrics.recordEDSAttempt(ctx, attempt, true)
			return eds, nil
		case errors.Is(getErr, context.DeadlineExceeded),
			errors.Is(getErr, context.Canceled):
		case errors.Is(getErr, p2p.ErrNotFound):
			getErr = share.ErrNotFound
			setStatus(peers.ResultCooldownPeer)
		case errors.Is(getErr, p2p.ErrInvalidResponse):
			setStatus(peers.ResultBlacklistPeer)
		default:
			setStatus(peers.ResultCooldownPeer)
		}

		if !ErrorContains(err, getErr) {
			err = errors.Join(err, getErr)
		}
		log.Debugw("eds: request failed",
			"hash", root.String(),
			"peer", peer.String(),
			"attempt", attempt,
			"err", getErr,
			"finished (s)", time.Since(reqStart))
	}
}

func (sg *ShrexGetter) GetSharesByNamespace(
	ctx context.Context,
	root *share.Root,
	id namespace.ID,
) (share.NamespacedShares, error) {
	var (
		attempt int
		err     error
	)
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		attempt++
		start := time.Now()
		peer, setStatus, getErr := sg.peerManager.Peer(ctx, root.Hash())
		if getErr != nil {
			err = errors.Join(err, getErr)
			log.Debugw("nd: couldn't find peer",
				"hash", root.String(),
				"err", getErr,
				"finished (s)", time.Since(start))
			sg.metrics.recordNDAttempt(ctx, attempt, false)
			return nil, fmt.Errorf("getter/shrex: %w", err)
		}

		reqStart := time.Now()
		reqCtx, cancel := ctxWithSplitTimeout(ctx, sg.minAttemptsCount-attempt+1, sg.minRequestTimeout)
		nd, getErr := sg.ndClient.RequestND(reqCtx, root, id, peer)
		cancel()
		switch {
		case getErr == nil:
			setStatus(peers.ResultNoop)
			sg.metrics.recordNDAttempt(ctx, attempt, true)
			return nd, nil
		case errors.Is(getErr, context.DeadlineExceeded),
			errors.Is(getErr, context.Canceled):
		case errors.Is(getErr, p2p.ErrNotFound):
			getErr = share.ErrNotFound
			setStatus(peers.ResultCooldownPeer)
		case errors.Is(getErr, p2p.ErrInvalidResponse):
			setStatus(peers.ResultBlacklistPeer)
		default:
			setStatus(peers.ResultCooldownPeer)
		}

		if !ErrorContains(err, getErr) {
			err = errors.Join(err, getErr)
		}
		log.Debugw("nd: request failed",
			"hash", root.String(),
			"peer", peer.String(),
			"attempt", attempt,
			"err", getErr,
			"finished (s)", time.Since(reqStart))
	}
}
