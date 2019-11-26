package main

import (
	"fmt"
	"log"
	"os"

	"flashnews/engine"
)

var logger *log.Logger

func main() {
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

	en.Run()
	fmt.Scan()
	os.Exit(0)
}
