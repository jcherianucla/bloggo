package main

import (
	"github.com/jcherianucla/bloggo/clients"
	"github.com/jcherianucla/bloggo/config"
	"github.com/jcherianucla/bloggo/handlers"
	"github.com/jcherianucla/bloggo/handlers/prometheus"
	"go.uber.org/fx"
)

func main() {
	fx.New(fxOptions()...).Run()
}

func fxOptions() []fx.Option {
	return []fx.Option{
		config.Module,
		prometheus.Module,
		clients.Module,
		handlers.Module,
	}
}
