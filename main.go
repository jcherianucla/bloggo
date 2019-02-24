package main

import (
	"github.com/jcherianucla/bloggo/clients"
	"github.com/jcherianucla/bloggo/config"
	"go.uber.org/fx"
)

func main() {
	fx.New(fxOptions()...).Run()
}

func fxOptions() []fx.Option {
	return []fx.Option{
		config.Module,
		clients.Module,
	}
}
