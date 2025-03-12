//go:build android
// +build android

package storage

import "path/filepath"

const appDirectory = "/data/data/com.groupalpha.agentofthings"

func getAppCacheDir() (string, error) {
	return filepath.Join(appDirectory, "cache"), nil
}

func getAppConfigDir() (string, error) {
	return filepath.Join(appDirectory, "files"), nil
}
