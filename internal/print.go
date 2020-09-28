package internal

import (
	"fmt"
	"os"
)

func ErrorOutputExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return
}
