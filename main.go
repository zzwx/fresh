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
	configPath := flag.String("c", "", fmt.Sprintf("config file path. Default is %q", runner.DefaultConfigPath))
	generate := flag.Bool("g", false, fmt.Sprintf("generate a sample settings file either at %q or at specified by -c location", runner.DefaultConfigPath))
	env := flag.String("e", "", fmt.Sprintf("environment variables prefix. %q is a default prefix", runner.EnvPrefix))

	flag.Parse()
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
