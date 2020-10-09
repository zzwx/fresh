package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func build() (string, bool) {
	args := []string{
		"go", "build", "-o", buildPath(),
	}
	args = append(args, strings.Fields(buildArgs())...)
	args = append(args, mainPath())
	buildLog("Building... %v", args)

	if mustUseDelve() {
		return "", true
	}

	cmd := exec.Command("go", args[1:]...) // [1: skips the "go" in args list

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
