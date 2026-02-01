package commands

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
)

type ShellCommand struct{}

func (c ShellCommand) Execute(args []string) (io.ReaderAt, uint64) {
	var out []byte
	var err error

	cmdStr := strings.Join(args, " ")
	out, err = exec.Command("cmd.exe", "/C", cmdStr).Output()
	if err != nil {
		out = []byte(err.Error())
	}
	return bytes.NewReader(out), uint64(len(out))
}
