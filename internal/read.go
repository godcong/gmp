package internal

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type BackupCallback func([]byte) ([]byte, error)

var modName = "go.mod"
var backupName = modName + ".old"
var ErrFileIsAlreadyExist = errors.New("target file is already exist")

func UpdateGoFile(path string, src, target string) error {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, info := range dir {
		if info.IsDir() {
			if err := UpdateGoFile(filepath.Join(path, info.Name()), src, target); err != nil {
				return err
			}
			continue
		}
		if filepath.Ext(info.Name()) == ".go" {
			if err := FileBackupFix(path, info.Name(), info.Name()+".old", nil); err != nil {
				return err
			}

			if err := fixImport(filepath.Join(path, info.Name()), src, target); err != nil {
				return err
			}
			continue
		}
		//skip
	}
	return nil
}

func fixImport(path string, srcModule string, targetModule string) error {
	//todo
	return nil
}

func FileBackupFix(path, srcName, targetName string, cb BackupCallback) error {
	if err := backupFile(path, srcName, targetName); err != nil {
		return err
	}

	open, err := os.Open(filepath.Join(path, targetName))
	if err != nil {
		return err
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	file, err := os.OpenFile(filepath.Join(path, srcName), os.O_TRUNC|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if cb != nil {
			line, err = cb(line)
			if err != nil {
				return err
			}
		}
		if _, err := file.Write(line); err != nil {
			return err
		}
		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}
	return nil
}

func getModule(src []byte) (string, bool) {
	if strings.Index(string(src), "module") >= 0 {
		if s := strings.Split(string(src), " "); len(s) > 1 {
			return s[len(s)-1], true
		}
	}
	return "src", false
}

func fixModule(src []byte, target string) []byte {
	if strings.Index(string(src), "module") >= 0 {
		return []byte("module " + target)
	}
	return src
}

func backupFile(path string, src, target string) error {
	back := filepath.Join(path, target)
	_, err := os.Stat(back)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	if err == nil {
		return ErrFileIsAlreadyExist
	}

	if err := os.Rename(filepath.Join(path, src), back); err != nil {
		return err
	}
	return nil
}
