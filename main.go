/*
Fresh is a command-line tool that builds, starts and restarts your written in Go application, including a web app, every time you save a Go or template file or any desired files you specify via configuration file.
Fresh can be used even without configuration file, using default values.

Fresh watches for file events, and every time a file is created, deleted or modified it will build and restart the application.

If `go build` returns an error, it will create a log file in the tmp folder and keep watching, attempting to rebuild if initial compilation was successful.
It will also attempt to kill previously created processes.

This is a fork of an original fresh (https://github.com/pilu/fresh) that is announced as unmaintained.

*/
package main

import (
	"flag"
	"fmt"
	"github.com/zzwx/fresh/runner"
	"os"
)

const VERSION = "1.3.2"

func main() {
	var version bool
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "prints current version and exits")
	configPath := flag.String("c", "", fmt.Sprintf("config file path. Default is %q", runner.DefaultConfigPath))
	generate := flag.Bool("g", false, fmt.Sprintf("generate a sample settings file either at %q or at specified by -c location", runner.DefaultConfigPath))
	env := flag.String("e", "", fmt.Sprintf("environment variables prefix. %q is a default prefix", runner.EnvPrefix))

	flag.Parse()
	if version || (len(os.Args) > 1 && os.Args[1] == "version") {
		fmt.Println(VERSION)
		return
	}

	if *env != "" {
		runner.EnvPrefix = *env
		fmt.Printf("Environment variables prefix set to %q\n", runner.EnvPrefix)
	}

	if *configPath != "" {
		if _, err := os.Stat(*configPath); err != nil {
			if *generate {
				runner.SaveRunnerConfigSettings(*configPath)
			} else {
				fmt.Printf("Can't find config file %q\n", *configPath)
				os.Exit(1)
			}
		} else {
			os.Setenv(runner.EnvPrefix+"CONFIG_PATH", *configPath) // RUNNER_CONFIG_PATH
		}
	} else {
		if *generate {
			runner.SaveRunnerConfigSettings(runner.DefaultConfigPath)
		}
	}

	if !*generate {
		runner.Start()
	}
}
