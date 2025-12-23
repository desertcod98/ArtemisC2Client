package commands

import (
    "os/exec"
)

type WhoamiCommand struct {}

func (c WhoamiCommand) Execute(args []string) <-chan string {
    result := make(chan string)
    go func() {
        out, err := exec.Command("whoami").Output()
        if err != nil {
            result <- err.Error()
        } else {
            result <- string(out)
        }
        close(result)
    }()
    return result
}
