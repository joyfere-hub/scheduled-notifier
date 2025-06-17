package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
	"github.com/joyfere-hub/scheduled-notifier/internal/ctx"
	"github.com/joyfere-hub/scheduled-notifier/internal/work"
	"github.com/segmentfault/pacman/log"
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
	systray.SetIcon(icon.Data)
	systray.SetTitle("scheduled-notifier")
	systray.SetTooltip("running")
	go func() {
		menuItem := systray.AddMenuItem("Quit", "Quit the whole app")
		<-menuItem.ClickedCh
		fmt.Println("Quit")
		app.exit()
	}()

	w, err := work.NewWorker(app.ctx)
	if err != nil {
		log.Error(err)
		return
	}
	app.w = w
}

func (app *App) exit() {
	app.w.Close()
	os.Exit(0)
}
