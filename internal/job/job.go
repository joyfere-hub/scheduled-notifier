package job

import (
	"fmt"
	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
	"github.com/joyfere-hub/scheduled-notifier/notifier"
)

type JobClient interface {
	FetchMessages() (*notifier.Messages, error)
}

var (
	clients = make(map[string]func(*conf.JobConfig) (JobClient, error))
)

func NewClient(clientType string, ctx *conf.JobConfig) (JobClient, error) {
	clientsFunc, ok := clients[clientType]
	if !ok {
		return nil, fmt.Errorf("client type %s not found", clientType)
	}
	return clientsFunc(ctx)
}

func register(clientType string, newClientFunc func(*conf.JobConfig) (JobClient, error)) {
	clients[clientType] = newClientFunc
}
