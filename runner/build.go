package runner

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func build() (string, bool) {
	if mustUseDelve() {
		return "", true
	}
	parts := []string{"build", "-o"}
	if buildPath() != "" {
		parts = append(parts, buildPath())
	}
	if buildArgs() != "" {
		parts = append(parts, buildArgs())
	}
	if mainPath() != "" {
		parts = append(parts, mainPath())
	}
	cmd := Cmd("go", strings.Join(parts, " "))
	buildLog(cmd.SysProcAttr.CmdLine)

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

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}
