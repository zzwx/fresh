//go:build !windows
// +build !windows

package runner

import (
	"os/exec"
	"strings"
)

// Cmd constructs a raw exec.Cmd to let it parse arguments
// as if the came in from the command line
func Cmd(cmdName string, args string) *exec.Cmd {
	// Let the args be parsed by the exec.Command instead of strings.Fields
	// that splits them into separate exec.Comman args
	// TODO(zzwx): Might need to deal with quoting
	cmd := exec.Command(cmdName, strings.Fields(args)...)
	return cmd
}

func CmdStr(cmd *exec.Cmd) string {
	return cmd.String()
}
