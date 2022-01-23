package runner

import (
	"io"
	"os/exec"
	"syscall"
)

func run() {
	var cmd *exec.Cmd
	if mustUseDelve() {
		cmd = Cmd("dlv", delveArgs())
	} else {
		cmd = Cmd(buildPath(), runArgs())
	}
	runnerLog("Starts %v", cmd.SysProcAttr.CmdLine)

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

	if isDebug() {
		runnerLog("PID %d", cmd.Process.Pid)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		defer func() {
			killDoneChannel <- struct{}{}
		}()
		<-killChannel

		pid := cmd.Process.Pid
		runnerLog("Kills PID %d", pid)

		if err := cmd.Process.Kill(); err != nil {
			if isDebug() {
				runnerLog("Killing PID %d error: %v", pid, err)
			}
		}

		if exiting {
			resetTermColors()
			done <- struct{}{}
		}

		_, err := cmd.Process.Wait()
		if isDebug() {
			if err != nil {
				runnerLog("Exited PID %d with error: %v", pid, err)
			}
		}
	}()
}

// Cmd constructs a raw exec.Cmd to let it parse arguments
// as if the came in from the command line
func Cmd(cmdName string, args string) *exec.Cmd {
	// Let the args be parsed by the exec.Command instead of strings.Fields
	// that splits them into separate exec.Comman args
	// TODO(zzwx): Might need to deal with quoting
	cmd := exec.Command(cmdName) // , strings.Fields(args)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.CmdLine = syscall.EscapeArg(cmd.Path) + " " + args
	return cmd
}
