package commands

import (
	"bytes"
	"io"
	"os/exec"
)

type ShellCommand struct{}

func (c ShellCommand) Execute(args []string) (io.ReaderAt, uint64) {
    var out []byte
    var err error
    out, err = exec.Command(args[0], args[1:]...).Output()
    if err != nil {
        out = []byte(err.Error())
    }
    return bytes.NewReader(out), uint64(len(out))
}
