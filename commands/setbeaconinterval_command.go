package commands

import "strconv"

type SetBeaconIntervalCommand struct {
	SetBeaconIntervalCh chan (int)
}

func (c SetBeaconIntervalCommand) Execute(args []string) []byte {
	if len(args) != 1 {
		return []byte("err: missing arg <seconds>")
	}
	interval, err := strconv.Atoi(args[0])
	if err != nil {
		return []byte("err: <seconds> must be int")
	}
	c.SetBeaconIntervalCh <- interval
	return []byte("ok")
}
