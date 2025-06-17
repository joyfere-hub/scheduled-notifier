package work

import (
	"fmt"
	"github.com/joyfere-hub/scheduled-notifier/internal/ctx"
	"github.com/joyfere-hub/scheduled-notifier/internal/job/private"
	"github.com/joyfere-hub/scheduled-notifier/notifier"
	"github.com/robfig/cron/v3"
	"github.com/segmentfault/pacman/log"
)

type Worker struct {
	cron *cron.Cron
}

func NewWorker(ctx *ctx.Context) (*Worker, error) {
	c := cron.New()
	messageChan := make(chan *[]notifier.Message)
	for _, jobConfig := range *ctx.Conf.Jobs {
		client, err := private.NewClient(&jobConfig)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Check client, fetch messages.\n")
		_, err = client.FetchMessages()
		if err != nil {
			return nil, err
		}
		_, err = c.AddFunc(jobConfig.Interval, func() {
			fmt.Printf("Fetch messages begin.\n")
			messages, err := client.FetchMessages()
			fmt.Printf("Fetch messages end. count : %d\n", len(*messages))
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
					for _, message := range *messages {
						fmt.Printf("Send message: %v\n", message)
						if err := message.Send(); err != nil {
							log.Error(err)
						}
					}
				}
			}
		}
	}()

	return &Worker{cron: c}, nil
}

func (w *Worker) Close() {
	w.cron.Stop()
}
