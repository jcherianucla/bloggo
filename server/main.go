package main

import (
	"github.com/jcherianucla/bloggo/clients/instrumenter"
	"github.com/jcherianucla/bloggo/config"
	"github.com/jcherianucla/bloggo/handlers"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func start(config config.AppConfig, instrumenter instrumenter.Instrument, gs *grpc.Server) {
	lis, err := net.Listen("tcp", config.Config().GRPCConfig.HostPort)
	logger := instrumenter.Logger("debug")
	if err != nil {
		logger.Fatal("Failed to retrieve server port", zap.Error(err))
	}
	go func() {
		if err = gs.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()
}

func main() {
	fx.New(fxOptions()...).Run()
}

func fxOptions() []fx.Option {
	return []fx.Option{
		config.Module,
		instrumenter.Module,
		handlers.Module,
		fx.Invoke(start),
	}
}
