package commands

var Dispatcher = map[string]Command{
	"whoami": WhoamiCommand{},
}