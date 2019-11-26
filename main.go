package main

import (
	"log"
	"os"

	_ "flashnews/crawlers"
	"flashnews/engine"
	_ "flashnews/utils"
)

var logger *log.Logger

func main() {
	// Logger
	logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

	// Engine
	en := engine.Engine{}
	err := en.Init(logger, "./config.json")
	if err != nil {
		logger.Println("[INIT ERROR]", err)
		os.Exit(1)
	}

	en.Run()
}
