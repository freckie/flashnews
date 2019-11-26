package main

import (
	"fmt"
	"log"
	"os"

	"runtime"

	"flashnews/engine"
)

var logger *log.Logger

func main() {
	// Parallel Processing
	runtime.GOMAXPROCS(4)

	// Logger
	logger = log.New(os.Stdout, "LOG ", log.LstdFlags)

	// Engine
	en := engine.Engine{}
	err := en.Init(logger, "config.json")
	if err != nil {
		logger.Println("[INIT ERROR]", err)
		fmt.Scan()
		os.Exit(0)
	}
	en.TG.TestMessage()

	en.Run()
	fmt.Scan()
	os.Exit(0)
}
