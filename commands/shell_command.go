package commands

import (
	"os/exec"
)

type ShellCommand struct{}

func (c ShellCommand) Execute(args []string) []byte {
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return []byte(err.Error())
	} else {
		return out
	}
}
