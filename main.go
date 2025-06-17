package main

import (
	"flag"
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
	"github.com/joyfere-hub/scheduled-notifier/internal/ctx"
	"github.com/joyfere-hub/scheduled-notifier/internal/work"
	"github.com/joyfere-hub/scheduled-notifier/res"
)

type App struct {
	ctx *ctx.Context
	w   *work.Worker
}

func main() {
	configPath := ""
	flag.StringVar(&configPath, "config", "", "config file path")
	flag.Parse()

	c, err := conf.ReadConfig(configPath)
	if err != nil {
		panic(err)
	}

	app := App{
		ctx: &ctx.Context{
			Conf: c,
		},
	}

	systray.Run(app.run, app.exit)
}

func (app *App) run() {
	systray.SetIcon(res.Icon)
	systray.SetTitle("")
	systray.SetTooltip("running")
	fetchItem := systray.AddMenuItem("Fetch", "Fetch messages")
	menuItem := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-fetchItem.ClickedCh:
				app.w.Fetch()
			case <-menuItem.ClickedCh:
				app.exit()
				return
			}
		}
	}()

	w, err := work.NewWorker(app.ctx)
	if err != nil {
		log.Panic(err)
		return
	}
	app.w = w
}

func (app *App) exit() {
	app.w.Close()
	os.Exit(0)
}
