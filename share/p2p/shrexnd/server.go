package shrexnd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/minio/sha256-simd"
	"go.uber.org/zap"

	"github.com/furyaxyz/go-libp2p-messenger/serde"

	"github.com/furyaxyz/elysium-node/share"
	"github.com/furyaxyz/elysium-node/share/eds"
	"github.com/furyaxyz/elysium-node/share/ipld"
	"github.com/furyaxyz/elysium-node/share/p2p"
	pb "github.com/furyaxyz/elysium-node/share/p2p/shrexnd/pb"
)

// Server implements server side of shrex/nd protocol to serve namespaced share to remote
// peers.
type Server struct {
	cancel context.CancelFunc

	host       host.Host
	protocolID protocol.ID

	getter share.Getter
	store  *eds.Store

	params     *Parameters
	middleware *p2p.Middleware
	metrics    *p2p.Metrics
}

// NewServer creates new Server
func NewServer(params *Parameters, host host.Host, store *eds.Store, getter share.Getter) (*Server, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("shrex-nd: server creation failed: %w", err)
	}

	srv := &Server{
		getter:     getter,
		store:      store,
		host:       host,
		params:     params,
		protocolID: p2p.ProtocolID(params.NetworkID(), protocolString),
		middleware: p2p.NewMiddleware(params.ConcurrencyLimit),
	}

	return srv, nil
}

// Start starts the server
func (srv *Server) Start(context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	srv.cancel = cancel

	handler := func(s network.Stream) {
		srv.handleNamespacedData(ctx, s)
	}
	srv.host.SetStreamHandler(srv.protocolID, srv.middleware.RateLimitHandler(handler))
	return nil
}

// Stop stops the server
func (srv *Server) Stop(context.Context) error {
	srv.cancel()
	srv.host.RemoveStreamHandler(srv.protocolID)
	return nil
}

func (srv *Server) observeRateLimitedRequests() {
	numRateLimited := srv.middleware.DrainCounter()
	if numRateLimited > 0 {
		srv.metrics.ObserveRequests(context.Background(), numRateLimited, p2p.StatusRateLimited)
	}
}

func (srv *Server) handleNamespacedData(ctx context.Context, stream network.Stream) {
	logger := log.With("peer", stream.Conn().RemotePeer())
	logger.Debug("server: handling nd request")

	srv.observeRateLimitedRequests()

	err := stream.SetReadDeadline(time.Now().Add(srv.params.ServerReadTimeout))
	if err != nil {
		logger.Debugw("server: setting read deadline", "err", err)
	}

	var req pb.GetSharesByNamespaceRequest
	_, err = serde.Read(stream, &req)
	if err != nil {
		logger.Warnw("server: reading request", "err", err)
		stream.Reset() //nolint:errcheck
		return
	}
	logger = logger.With("namespaceId", string(req.NamespaceId), "hash", string(req.RootHash))
	logger.Debugw("server: new request")

	err = stream.CloseRead()
	if err != nil {
		logger.Debugw("server: closing read side of the stream", "err", err)
	}

	err = validateRequest(req)
	if err != nil {
		logger.Debugw("server: invalid request", "err", err)
		stream.Reset() //nolint:errcheck
		return
	}

	ctx, cancel := context.WithTimeout(ctx, srv.params.HandleRequestTimeout)
	defer cancel()

	dah, err := srv.store.GetDAH(ctx, req.RootHash)
	if err != nil {
		if errors.Is(err, eds.ErrNotFound) {
			srv.respondNotFoundError(ctx, logger, stream)
			return
		}
		logger.Errorw("server: retrieving DAH", "err", err)
		srv.respondInternalError(ctx, logger, stream)
		return
	}

	shares, err := srv.getter.GetSharesByNamespace(ctx, dah, req.NamespaceId)
	if errors.Is(err, share.ErrNotFound) {
		srv.respondNotFoundError(ctx, logger, stream)
		return
	}
	if err != nil {
		logger.Errorw("server: retrieving shares", "err", err)
		srv.respondInternalError(ctx, logger, stream)
		return
	}

	resp := namespacedSharesToResponse(shares)
	srv.respond(ctx, logger, stream, resp)
}

// validateRequest checks correctness of the request
func validateRequest(req pb.GetSharesByNamespaceRequest) error {
	if len(req.NamespaceId) != ipld.NamespaceSize {
		return fmt.Errorf("incorrect namespace id length: %v", len(req.NamespaceId))
	}
	if len(req.RootHash) != sha256.Size {
		return fmt.Errorf("incorrect root hash length: %v", len(req.RootHash))
	}

	return nil
}

// respondNotFoundError sends internal error response to client
func (srv *Server) respondNotFoundError(ctx context.Context,
	logger *zap.SugaredLogger, stream network.Stream) {
	resp := &pb.GetSharesByNamespaceResponse{
		Status: pb.StatusCode_NOT_FOUND,
	}
	srv.respond(ctx, logger, stream, resp)
}

// respondInternalError sends internal error response to client
func (srv *Server) respondInternalError(ctx context.Context,
	logger *zap.SugaredLogger, stream network.Stream) {
	resp := &pb.GetSharesByNamespaceResponse{
		Status: pb.StatusCode_INTERNAL,
	}
	srv.respond(ctx, logger, stream, resp)
}

// namespacedSharesToResponse encodes shares into proto and sends it to client with OK status code
func namespacedSharesToResponse(shares share.NamespacedShares) *pb.GetSharesByNamespaceResponse {
	rows := make([]*pb.Row, 0, len(shares))
	for _, row := range shares {
		proof := &pb.Proof{
			Start: int64(row.Proof.Start()),
			End:   int64(row.Proof.End()),
			Nodes: row.Proof.Nodes(),
		}

		row := &pb.Row{
			Shares: row.Shares,
			Proof:  proof,
		}

		rows = append(rows, row)
	}

	return &pb.GetSharesByNamespaceResponse{
		Status: pb.StatusCode_OK,
		Rows:   rows,
	}
}

func (srv *Server) respond(ctx context.Context,
	logger *zap.SugaredLogger, stream network.Stream, resp *pb.GetSharesByNamespaceResponse) {
	err := stream.SetWriteDeadline(time.Now().Add(srv.params.ServerWriteTimeout))
	if err != nil {
		logger.Debugw("server: setting write deadline", "err", err)
	}

	_, err = serde.Write(stream, resp)
	if err != nil {
		logger.Warnw("server: writing response", "err", err)
		stream.Reset() //nolint:errcheck
		return
	}

	switch {
	case resp.Status == pb.StatusCode_OK:
		srv.metrics.ObserveRequests(ctx, 1, p2p.StatusSuccess)
	case resp.Status == pb.StatusCode_NOT_FOUND:
		srv.metrics.ObserveRequests(ctx, 1, p2p.StatusNotFound)
	case resp.Status == pb.StatusCode_INTERNAL:
		srv.metrics.ObserveRequests(ctx, 1, p2p.StatusInternalErr)
	}
	if err = stream.Close(); err != nil {
		logger.Debugw("server: closing stream", "err", err)
	}
}
