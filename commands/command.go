package commands

type Command interface{
	Execute(args []string) <- chan string
}