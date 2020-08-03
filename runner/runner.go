package runner

import (
	"io"
	"os/exec"
	"strings"
)

func run() bool {
	if isDebug() {
		runnerLog("Running...")
	}

	var cmd *exec.Cmd

	if mustUseDelve() {
		cmd = exec.Command("dlv", strings.Fields(delveArgs())...)
	} else {
		cmd = exec.Command(buildPath(), strings.Fields(buildArgs())...)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		if isDebug() {
			runnerLog("Stopping...")
		}

		pid := cmd.Process.Pid
		runnerLog("Killing PID %d...", pid)

		if err := cmd.Process.Kill(); err != nil {
			if isDebug() {
				runnerLog("Killing PID %d failed: %v", pid, err)
			}
		}

		if exiting == true {
			resetTermColors()
			done <- true
		}

		if isDebug() {
			runnerLog("Killed")
		}
		cmd.Process.Wait()

	}()

	return true
}
