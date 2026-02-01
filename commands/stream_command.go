package commands

import "io"

type StreamCommand interface {
	Execute(args []string) (io.ReaderAt, uint64)
}
