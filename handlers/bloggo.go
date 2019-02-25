package handlers

import (
	"context"
	"net"

	"go.uber.org/zap"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/jcherianucla/bloggo/clients/instrumenter"
	"github.com/jcherianucla/bloggo/config"

	"go.uber.org/fx"

	"github.com/jcherianucla/bloggo/.gen/idl/proto"
)

var Module = fx.Provide(New)

var panicHandler = grpc_recovery.RecoveryHandlerFunc(func(p interface{}) error {
	return status.Errorf(codes.Internal, "%s", p)
})

type Params struct {
	fx.In

	config.AppConfig
	instrumenter.Instrument
}

type bloggoHandler struct {
}

func New(p Params) {
	lis, err := net.Listen("tcp", p.Config().GRPCConfig.HostPort)
	logger := p.Logger("debug")
	if err != nil {
		logger.Fatal("Failed to retrieve server port", zap.Error(err))
	}
	//TODO: Add TLS and auth interceptor
	gs := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(panicHandler)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(panicHandler)),
		)))
	grpc_prometheus.Register(gs)
	proto.RegisterBloggoServer(gs, bloggoHandler{})
	go func() {
		if err = gs.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()
}

func (bh bloggoHandler) CreatePost(
	ctx context.Context,
	in *proto.CreatePostRequest,
) (*proto.CreatePostResponse, error) {
	return nil, nil
}
func (bh bloggoHandler) FetchPost(
	ctx context.Context,
	in *proto.FetchPostRequest,
) (*proto.FetchPostResponse, error) {
	return nil, nil
}
func (bh bloggoHandler) UpdatePost(
	ctx context.Context,
	in *proto.UpdatePostRequest,
) (*proto.UpdatePostResponse, error) {
	return nil, nil
}
func (bh bloggoHandler) DeletePost(
	ctx context.Context,
	in *proto.DeletePostRequest,
) (*proto.DeletePostResponse, error) {
	return nil, nil
}
