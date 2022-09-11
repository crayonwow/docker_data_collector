package internal

import (
	"docker_data_collector/internal/config"
	"docker_data_collector/internal/graceful"
	"docker_data_collector/internal/notifier"
	"docker_data_collector/internal/telegrambot"
	"docker_data_collector/internal/watcher"

	"github.com/sirupsen/logrus"
)

func Run() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	gCtx := graceful.Context()

	bot, err := telegrambot.NewBot(cfg.TG)
	if err != nil {
		panic(err)
	}

	w, err := watcher.NewWatcher()
	if err != nil {
		panic(err)
	}

	notifier.Prepare(gCtx, bot, w)

	logrus.Info("Started")

	// todo gracefully stop services
	<-gCtx.Done()
}
