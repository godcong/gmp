package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileBackup(t *testing.T) {
	type args struct {
		path   string
		src    string
		target string
		cb     BackupCallback
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				path:   "D:\\workspace\\golang\\project\\gmp\\test",
				src:    "go.mod",
				target: "go.mod.old",
				cb:     nil,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				path:   "D:\\workspace\\golang\\project\\gmp\\test",
				src:    "go.mod",
				target: "go.mod.old",
				cb: func(bytes []byte) []byte {
					str := string(bytes)
					if strings.Index(str, "module") >= 0 {
						return []byte("module " + "new.package")
					}
					return bytes
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FileBackup(tt.args.path, tt.args.src, tt.args.target, tt.args.cb)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			os.Remove(filepath.Join(tt.args.path, tt.args.target))
		})
	}
}
