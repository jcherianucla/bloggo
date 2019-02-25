package clients

import (
	"github.com/jcherianucla/bloggo/config"
	"github.com/jcherianucla/bloggo/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In
	config.AppConfig
}

type Result struct {
	fx.Out
	Instrument
}

type Instrument interface {
	Logger(logger string) *zap.Logger
	Close()
}

func New(p Params) Result {
	logStyle := p.Config().LoggerConfig.Style
	return Result{
		Instrument: &instrument{
			style: logStyle,
		},
	}
}

type instrument struct {
	style   string
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

func (i *instrument) Close() {
	for _, log := range i.loggers {
		err := log.Sync()
		utils.HandleErr(err)
	}
}
