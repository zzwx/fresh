package runner

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	watchChannel    chan string
	killChannel     chan struct{}
	killDoneChannel chan struct{}
	doneChannel     chan struct{}
	exiting         bool
	mainLog         logFunc
	watcherLog      logFunc
	runnerLog       logFunc
	buildLog        logFunc
	appLog          logFunc
)

func flushEvents() {
	for {
		select {
		case eventName := <-watchChannel:
			if isDebug() {
				mainLog("Skipping %s", eventName)
			}
		default:
			return
		}
	}
}

func start() {
	loopIndex := 0
	delay := buildDelay()

	started := false

	go termHandler()

	go func() {
		for {
			if isDebug() {
				mainLog("Waiting for (Loop: %d, Goroutines: %d)", loopIndex, runtime.NumGoroutine())
			}
			eventName := <-watchChannel

			if isDebug() {
				mainLog("First event: %s", eventName)
			}
			if loopIndex > 0 {
				if isDebug() {
					mainLog("Sleeping %v", delay)
				}
				time.Sleep(delay)
			}

			err := removeBuildErrorsLog()
			if err != nil && !os.IsNotExist(err) {
				if isDebug() {
					mainLog(err.Error())
				}
			}

			if isDebug() {
				mainLog("Skipping events")
			}
			flushEvents()

			buildFailed := false
			if shouldRebuild(eventName) {
				mainLog("Rebuilding due to: %v", eventName)
				err := build()
				if err != nil {
					buildFailed = true
					mainLog("Failed:\n%v", err)
					if !started {
						os.Exit(1)
					}
					createBuildErrorsLog(err.Error())
				}

				if !buildFailed {
					if started {
						killChannel <- struct{}{}
						<-killDoneChannel // we don't want run() before killing is finished
					}
					run()
					started = true
				}
			}
			loopIndex++
		}
	}()
}

func init() {
	watchChannel = make(chan string, 1000)
	killChannel = make(chan struct{})
	killDoneChannel = make(chan struct{})
}

const maxPrefixLength = 7

func initLogFuncs() {
	mainLog = newLogFunc("main", true)
	watcherLog = newLogFunc("watcher", true)
	runnerLog = newLogFunc("runner", true)
	buildLog = newLogFunc("build", true)
	appLog = newLogFunc("app", false)
}

func setEnvVars() {
	os.Setenv("DEV_RUNNER", "1")
	wd, err := os.Getwd()
	if err == nil {
		os.Setenv(EnvPrefix+"WD", wd) // RUNNER_WD
	}
	os.Setenv(EnvPrefix+"CONFIG_PATH", ConfigPath)
	t := reflect.TypeOf(Settings{})
	v := reflect.ValueOf(settings)
	for i := 0; i < t.NumField(); i++ {
		keyName, _, noenv := tagDetails(t.Field(i).Tag)
		if keyName != "" && !noenv {
			envKey := fmt.Sprintf("%s%s", EnvPrefix, strings.ToUpper(keyName))
			field := v.Field(i)
			err = nil
			if field.Kind() == reflect.String {
				err = os.Setenv(envKey, field.String())
			} else if field.Kind() == reflect.Bool {
				err = os.Setenv(envKey, strconv.FormatBool(field.Bool()))
			} else if field.Kind() == reflect.Uint {
				err = os.Setenv(envKey, strconv.FormatUint(field.Uint(), 10))
			}
			if err != nil {
				fmt.Printf("Error setting %q to %v due to %v", envKey, field.String(), err)
			}
			//runnerLog("%v set to %q", envKey, os.Getenv(envKey))
		}
	}

}

func termHandler() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	exiting = true
	killChannel <- struct{}{}
}

// Start Watches for file changes in the root directory.
// After each file system event it builds and (re)starts the application.
func Start() {

	doneChannel = make(chan struct{})

	initLimit()
	initLogFuncs() // Initialize log functions with default settings for initSettings to use
	initSettings()
	initLogFuncs() // Repeat after reading config file
	initFolders()
	setEnvVars()
	watch()
	start()
	watchChannel <- string(filepath.Separator)

	<-doneChannel
}
