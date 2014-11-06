package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	boshapp "github.com/cloudfoundry/bosh-agent/app"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

const mainLogTag = "main"

func main() {
	logger := boshlog.NewLogger(boshlog.LevelDebug)
	defer logger.HandlePanic("Main")
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	logger.Debug(mainLogTag, "Starting agent")

	app := boshapp.New(logger)

	err := app.Setup(os.Args)
	if err != nil {
		logger.Error(mainLogTag, "App setup %s", err.Error())
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		logger.Error(mainLogTag, "App run %s", err.Error())
		os.Exit(1)
	}

}
