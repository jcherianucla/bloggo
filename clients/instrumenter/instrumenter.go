package instrumenter

import (
	"github.com/jcherianucla/bloggo/config"
	"github.com/jcherianucla/bloggo/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides the Instrumenter through Fx
var Module = fx.Provide(New)

// Params defines the input dependencies for the Instrumenter
type Params struct {
	fx.In
	config.AppConfig
}

// Result defines the output dependency that is the Instrumenter
type Result struct {
	fx.Out
	Instrument
}

// Instrument defines the interface to interact with any instrumentation
// across the application
type Instrument interface {
	Metrics
	// Logger returns a zap Logger based on the type of logger desired
	Logger(logger string) *zap.Logger
	// Close gracefully shutdown the instrumentation clients
	Close()
}

// Metrics defines the wrapper interfaces around the common Prometheus
// instrument types
type Metrics interface {
	// Counter gives back a prometheus counter configured with opts
	Counter(opts prometheus.CounterOpts) prometheus.Counter
	// Histogram gives back a prometheus histogram configured with opts
	Histogram(opts prometheus.HistogramOpts) prometheus.Histogram
}

// New creates a new instrumenter
func New(p Params) Result {
	return Result{
		Instrument: &instrument{
			loggers: make(map[string]*zap.Logger),
		},
	}
}

type instrument struct {
	loggers map[string]*zap.Logger
}

func (i *instrument) Logger(logger string) *zap.Logger {
	switch logger {
	case utils.DebugLogType:
		if log, ok := i.loggers[utils.DebugLogType]; ok {
			return log
		}
		log, err := zap.NewDevelopment()
		utils.HandleErr(err)
		i.loggers[utils.DebugLogType] = log
		return log
	case utils.ProdLogType:
		fallthrough
	default:
		if log, ok := i.loggers[utils.ProdLogType]; ok {
			return log
		}
		log, err := zap.NewProduction()
		utils.HandleErr(err)
		i.loggers[utils.ProdLogType] = log
		return log
	}
}

func (i *instrument) Counter(opts prometheus.CounterOpts) prometheus.Counter {
	return promauto.NewCounter(opts)
}

func (i *instrument) Histogram(opts prometheus.HistogramOpts) prometheus.Histogram {
	return promauto.NewHistogram(opts)
}

func (i *instrument) Close() {
	for _, log := range i.loggers {
		err := log.Sync()
		utils.HandleErr(err)
	}
}
