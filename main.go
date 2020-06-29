//go:generate goversioninfo -file-version="v1.9.2"

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

		// Get config file path from os.Args
		var configPath string
		args := os.Args
		if len(args) <= 1 {
			configPath = "config.json"
		} else {
			configPath = args[1]
		}

		// Engine
		en := engine.Engine{}
		err := en.Init(logger, configPath)
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
