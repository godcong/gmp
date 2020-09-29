package internal

import "testing"

func TestFileBackup(t *testing.T) {
	type args struct {
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
				src:    "D:\\workspace\\golang\\project\\gmp\\test\\go.mod",
				target: "D:\\workspace\\golang\\project\\gmp\\test\\go.mod.old",
				cb:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FileBackup(tt.args.src, tt.args.target, tt.args.cb)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
