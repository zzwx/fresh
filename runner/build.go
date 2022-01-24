package runner

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func build() error {
	if mustUseDelve() {
		return nil
	}
	parts := []string{"build"}
	if buildPath() != "" {
		parts = append(parts, "-o")
		parts = append(parts, buildPath())
	}
	if buildArgs() != "" {
		parts = append(parts, buildArgs())
	}
	if mainPath() != "" {
		parts = append(parts, mainPath())
	}
	cmd := Cmd("go", strings.Join(parts, " "))
	buildLog("Building %v", CmdStr(cmd))

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
		return errors.New(string(errBuf))
	}

	return nil
}
