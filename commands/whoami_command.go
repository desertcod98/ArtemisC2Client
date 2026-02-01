package commands

import (
	"os/exec"
)

type WhoamiCommand struct{}

func (c WhoamiCommand) Execute(args []string) string {
	out, err := exec.Command("whoami").Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}
