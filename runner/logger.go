package runner

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
)

type logFunc func(string, ...interface{})

var logger = log.New(colorable.NewColorableStderr(), "", 0)

func resetTermColors() {
	if settings.Colors {
		logger.Printf(fmt.Sprintf("\033[%sm", colors["reset"]))
	}
}

func newLogFunc(prefix string, withEscape bool) func(string, ...interface{}) {
	var color, white, reset string
	if settings.Colors {
		lColor := logColor(prefix)
		if lColor != "" {
			color = fmt.Sprintf("\033[%sm", lColor)
		}
		white = fmt.Sprintf("\033[%sm", colors["white"])
		reset = fmt.Sprintf("\033[%sm", colors["reset"])
	}
	prefix = fmt.Sprintf("%-"+strconv.FormatInt(maxPrefixLength, 10)+"s", prefix)

	return func(msg string, v ...interface{}) {
		now := time.Now()
		timeString := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
		if withEscape {
			logger.Printf(fmt.Sprintf("%s%s %s |%s %s%s", color, timeString, prefix, white, msg, reset), v...)
		} else {
			// Message can arrive with new lines, possibly due to buffer flushing
			split := strings.Split(msg, "\n")
			for i, line := range split {
				if i == len(split)-1 {
					if len(line) == 0 { // Skip final \n that has nothing after it
						break
					}
				}
				logger.Print(fmt.Sprintf("%s%s %s |%s %s%s", color, timeString, prefix, white, line, reset))
			}
		}
	}
}

func fatal(err error) {
	logger.Fatal(err)
}

type appLogWriter struct{}

func (a appLogWriter) Write(p []byte) (n int, err error) {
	appLog(string(p))
	return len(p), nil
}
