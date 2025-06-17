package work

import (
	"log"

	"github.com/joyfere-hub/scheduled-notifier/internal/ctx"
	"github.com/joyfere-hub/scheduled-notifier/internal/job"
	"github.com/joyfere-hub/scheduled-notifier/notifier"
	"github.com/robfig/cron/v3"
)

type Worker struct {
	cron *cron.Cron
}

func NewWorker(ctx *ctx.Context) (*Worker, error) {
	c := cron.New()
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
					log.Panic(err)
				}
			}
		}
	}()

	return &Worker{cron: c}, nil
}

func (w *Worker) Close() {
	w.cron.Stop()
}
