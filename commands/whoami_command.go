package commands

import (
	"os/exec"
)

type WhoamiCommand struct{}

func (c WhoamiCommand) Execute(args []string) []byte {
	out, err := exec.Command("whoami").Output()
	if err != nil {
		return []byte(err.Error())
	}
	return out
}
