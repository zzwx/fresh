package runner

import (
	"io"
	"os/exec"
	"syscall"
)

func run() bool {
	var cmd *exec.Cmd
	if mustUseDelve() {
		cmd = Cmd("dlv", delveArgs())
	} else {
		cmd = Cmd(buildPath(), runArgs())
	}
	if isDebug() {
		runnerLog(cmd.SysProcAttr.CmdLine)
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
