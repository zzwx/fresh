package runner

import (
	"io"
	"os/exec"
	"strings"
)

func run() bool {
	runnerLog("Running...")

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
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		if err := cmd.Process.Kill(); err != nil {
			panic(err)
		}
	}()
	cmd.Wait()
	return true
}
