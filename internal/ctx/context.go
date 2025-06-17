package ctx

import (
	"context"
	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
)

type Context struct {
	context.Context
	AppQuitChan chan any
	Conf        *conf.BaseConfig
}
