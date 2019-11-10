package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func Homedir() string {
	return filepath.Join(home(), ".cwm")
}

func WorkerKeyPath() string {
	return filepath.Join(workingDirectory(), "worker_key")
}

func TsaHostKeyPath() string {
	return filepath.Join(workingDirectory(), "tsa_host_key.pub")
}

func workingDirectory() string {
	dir, _ := os.Getwd()
	return dir
}

func home() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}

	return os.Getenv("HOME")
}

