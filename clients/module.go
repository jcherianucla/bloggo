package clients

import (
	"github.com/jcherianucla/bloggo/clients/datastore"
	"github.com/jcherianucla/bloggo/clients/instrumenter"
	"go.uber.org/fx"
)

// Module defines the Fxified client modules
var Module = fx.Options(
	datastore.Module,
	instrumenter.Module,
)
