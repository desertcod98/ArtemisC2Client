package commands

import (
	"bytes"
	"io"
	"os"
)

type DownloadCommand struct{}

func (c DownloadCommand) Execute(args []string) (io.ReaderAt, int64, io.Closer) {
	var out []byte
	var err error

	filepath := args[0]
	file, err := os.Open(filepath)
	if err != nil {
		out = []byte(err.Error())
		return bytes.NewReader(out), int64(len(out)), nil
	}

	info, err := file.Stat()
	if err != nil {
		out = []byte(err.Error())
		return bytes.NewReader(out), int64(len(out)), nil
	}
	size := info.Size()

	return file, size, file
}
