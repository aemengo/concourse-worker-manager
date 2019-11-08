package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func Homedir() string {
	return filepath.Join(home(), ".cwm")
}

func home() string {
	if runtime.GOOS == "windows" {
		//TODO: make sure to return a real HOME dir
		return ""
	}

	return os.Getenv("HOME")
}

