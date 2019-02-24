package clients

import (
	"context"

	"github.com/jcherianucla/bloggo/.gen/idl/proto"

	"github.com/jcherianucla/bloggo/config"

	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In

	config.AppConfig
	proto.BloggoServer
}

type Result struct {
	fx.Out
	GRPCServer
}

type GRPCServer interface {
	CreateServer(ctx context.Context) error
}

func New(p Params) Result {
	return Result{}
}

type grpcserver struct {
	address string
}

func (gs *grpcserver) CreateServer(ctx context.Context) error {
	return nil
}
