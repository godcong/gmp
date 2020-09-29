package internal

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type BackupCallback func([]byte) []byte

var modName = "go.mod"
var backupName = modName + ".old"
var ErrFileIsAlreadyExist = errors.New("target file is already exist")

func UpdateGoFile(path string, src, target string) error {
	current := GetCurrentPath()
	dir, err := ioutil.ReadDir(current)
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
			if err := backupFile(path, info.Name(), info.Name()+".old"); err != nil {
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

func FileBackup(src, target string, cb BackupCallback) (string, error) {
	open, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	file, err := os.OpenFile(target, os.O_TRUNC|os.O_CREATE|os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		return "", err
	}
	defer file.Close()
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if cb != nil {
			line = cb(line)
		}
		if _, err := file.Write(line); err != nil {
			return "", err
		}
	}
	return src, nil
}

func fixModule(target string) []byte {
	//if strings.Index(string(line), "module") >= 0 {
	//	if s := strings.Split(string(line), " "); len(s) > 1 {
	//		src = s[len(s)-1]
	//	}
	//	line = fixModule(target)
	//}
	return []byte("module " + target)
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
