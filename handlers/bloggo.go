package handlers

import (
	"context"
	"github.com/jcherianucla/bloggo/idl/proto"
	"github.com/jcherianucla/bloggo/idl/proto/data"
	"github.com/jcherianucla/bloggo/idl/proto/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/jcherianucla/bloggo/clients/instrumenter"
	"github.com/jcherianucla/bloggo/config"

	"go.uber.org/fx"
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

type Result struct {
	fx.Out

	Server *grpc.Server
}

func New(p Params) Result {
	logger := p.Logger("debug")
	//TODO: Add TLS and auth interceptor
	gs := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(panicHandler)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(panicHandler)),
		)))
	proto.RegisterBloggoServer(gs, bloggoHandler{})

	return Result{Server: gs}
}

func (bh bloggoHandler) Create(
	ctx context.Context,
	in *proto.CreatePostRequest,
) (*proto.CreatePostResponse, error) {
	return &proto.CreatePostResponse{
		Data: &models.Post{
			ProtoUuid: &data.PostId{
				Uuid: &data.UUID{
					Value: "hello",
				},
			},
			Title:       "Hello World",
			Description: "This is the first static response from this gRPC service",
		},
	}, nil
}
