/*
Fresh is a command-line tool that builds and (re)starts your written in Go application, including a web app, every time you save a Go or template file or any desired files you specify using configuration.

Fresh can be started even without any configuration files.
It will watch for file events, and every time you create / modify or delete a file it will build and restart the application.

If `go build` returns an error, it will create a log file in the tmp folder and keep watching, attempting to rebuild. It will also attempt to kill previously created processes.

This is a fork of an original fresh (https://github.com/pilu/fresh) that is set as unmaintained. Check the README.md for more details.
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zzwx/fresh/runner"
)

func main() {
	configPath := flag.String("c", "", "config file path")
	flag.Parse()

	if *configPath != "" {
		if _, err := os.Stat(*configPath); err != nil {
			fmt.Printf("Can't find config file %q\n", *configPath)
			os.Exit(1)
		} else {
			os.Setenv("RUNNER_CONFIG_PATH", *configPath)
		}
	}

	runner.Start()
}
