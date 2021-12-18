package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	qa "querycsv/app"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type App interface {
	ReadConfigFile(string)
	GetQueryString()
	Run(ctx context.Context)
	Timeout() time.Duration
	Done() chan bool
}

func main() {
	var app App
	log := logrus.New()
	f, err := os.OpenFile(".\\log.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(logrus.TextFormatter)

	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {
		fmt.Println(err)
	} else {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}
	log.Formatter = new(logrus.JSONFormatter)

	app = qa.NewApp(log)

	var configPath *string = flag.String("conf", ".\\config\\config.yaml", "Configuration file's path")
	flag.Parse()

	app.ReadConfigFile((*configPath))

	app.GetQueryString()

	ctx, cancel := context.WithTimeout(context.Background(), app.Timeout())
	go app.Run(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		case <-sigCh:
			log.Info("got SIGINT")
			cancel()
		case <-app.Done():
			break loop
		}

	}

}
