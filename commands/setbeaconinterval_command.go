package commands

import "strconv"

type SetBeaconIntervalCommand struct {
	SetBeaconIntervalCh chan (int)
}

func (c SetBeaconIntervalCommand) Execute(args []string) <-chan string {
	resultCh := make(chan string)
	go func() {
		if len(args) != 1 {
			resultCh <- "err: missing arg <seconds>"
			close(resultCh)
			return
		}
		interval, err := strconv.Atoi(args[0])
		if err != nil {
			resultCh <- "err: <seconds> must be int"
			close(resultCh)
			return
		}
		c.SetBeaconIntervalCh <- interval
		resultCh <- "ok"
		close(resultCh)
	}()
	return resultCh
}
