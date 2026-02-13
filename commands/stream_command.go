package commands

import "io"

type StreamCommand interface {
	Execute(args []string) (io.ReaderAt, int64, io.Closer)
}
