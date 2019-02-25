package handlers

import (
	"net/http"

	"github.com/jcherianucla/bloggo/utils"

	"github.com/jcherianucla/bloggo/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

const (
	_metricsEndpoint = "/metrics"
)

var PrometheusModule = fx.Provide(NewPrometheus)

type Params struct {
	fx.In
	config.AppConfig
}

func NewPrometheus(p Params) {
	http.Handle(_metricsEndpoint, promhttp.Handler())
	err := http.ListenAndServe(p.Config().MetricsConfig.Port, nil)
	if err != nil {
		utils.HandleErr(err)
	}
}
