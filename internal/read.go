package internal

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var modName = "go.mod"
var backupName = modName + ".old"
var ErrFileIsAlreadyExist = errors.New("target file is already exist")

func UpdateMod(target string) (string, error) {
	path := GetCurrentPath()
	back := filepath.Join(path, backupName)
	open, err := os.Open(back)
	if err != nil {
		return "", err
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	file, err := os.OpenFile(filepath.Join(path, modName), os.O_TRUNC|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var src string
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if strings.Index(string(line), "module") >= 0 {
			if s := strings.Split(string(line), " "); len(s) > 1 {
				src = s[len(s)-1]
			}
			line = fixModule(target)
		}
		if _, err := file.Write(line); err != nil {
			return "", err
		}
	}
	return src, nil
}

func fixModule(target string) []byte {
	return []byte("module " + target)
}

func backupMod(path string) error {
	back := filepath.Join(path, backupName)
	_, err := os.Stat(back)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	if err == nil {
		return ErrFileIsAlreadyExist
	}

	if err := os.Rename(filepath.Join(path, modName), back); err != nil {
		return err
	}
	return nil
}
