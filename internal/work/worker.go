package work

import (
	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
	"github.com/joyfere-hub/scheduled-notifier/internal/ctx"
	"github.com/joyfere-hub/scheduled-notifier/internal/job"
	"github.com/joyfere-hub/scheduled-notifier/notifier"
	"github.com/robfig/cron/v3"
	"log"
	"maps"
)

type Worker struct {
	cron         *cron.Cron
	jobClientMap *map[string]*job.JobClient
	jobConfigMap *map[string]*conf.JobConfig
	messageChan  chan *notifier.Messages
}

func NewWorker(ctx *ctx.Context) (*Worker, error) {
	c := cron.New()
	jobClientMap := make(map[string]*job.JobClient)
	jobConfigMap := make(map[string]*conf.JobConfig)
	messageChan := make(chan *notifier.Messages)
	for _, jobConfig := range *ctx.Conf.Jobs {
		client, err := job.NewClient(jobConfig.Type, &jobConfig)
		if err != nil {
			return nil, err
		}
		log.Printf("Check %s client, fetch messages.", jobConfig.Type)
		messages, err := client.FetchMessages()
		if err != nil {
			return nil, err
		}
		err = messages.Send()
		if err != nil {
			return nil, err
		}
		jobClientMap[jobConfig.Type] = &client
		jobConfigMap[jobConfig.Type] = &jobConfig
		_, err = c.AddFunc(jobConfig.Interval, func() {
			log.Printf("Fetch messages begin.")
			messages, err := client.FetchMessages()
			log.Printf("Fetch messages end. count : %d", len(*messages))
			if err != nil {
				panic(err)
			}
			messageChan <- messages
		})

		if err != nil {
			return nil, err
		}
	}

	c.Start()

	go func() {
		for {
			select {
			case messages := <-messageChan:
				if messages != nil {
					// TODO message group
					err := messages.Send()
					if err != nil {
						log.Printf("Send message fail, %v", err)
					}
				}
			}
		}
	}()

	return &Worker{cron: c, jobClientMap: &jobClientMap, jobConfigMap: &jobConfigMap, messageChan: messageChan}, nil
}
func (w *Worker) Fetch() {
	clients := maps.Values(*w.jobClientMap)
	for client := range clients {
		log.Printf("Fetch messages begin.")
		messages, err := (*client).FetchMessages()
		log.Printf("Fetch messages end. count : %d.", len(*messages))
		if err != nil {
			continue
		}
		w.messageChan <- messages
	}
}

func (w *Worker) Close() {
	w.cron.Stop()
}
