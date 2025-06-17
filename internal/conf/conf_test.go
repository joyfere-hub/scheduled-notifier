package conf

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestReadConfig(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(currentFile), "test")

	type args struct {
		configFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantC   *BaseConfig
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				configFilePath: filepath.Join(basePath, "base.yaml"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := ReadConfig(tt.args.configFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf(gotC.String())
		})
	}
}
