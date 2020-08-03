package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func build() (string, bool) {
	buildLog("Building...")

	if mustUseDelve() {
		return "", true
	}

	cmd := exec.Command("go", "build", "-o", buildPath(), filepath.Join(root(), mainPath()))

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
