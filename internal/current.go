package internal

import "os"

func GetCurrentPath() string {
	wd, err := os.Getwd()
	ErrorOutputExit(err)
	return wd
}
