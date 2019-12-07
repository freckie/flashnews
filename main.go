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
	// Testmode
	testMode := false

	if !testMode {
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

		// Parallel Processing
		runtime.GOMAXPROCS(en.Cfg.Crawler.MaxProcs)

		en.Run()
		fmt.Scan()
		os.Exit(0)
	} else {

		/* Write Test Code Here */

	}
}
