package notifier

import (
	"runtime"
	"testing"
)

func Test_sendNotification(t *testing.T) {
	type args struct {
		title   string
		message string
		link    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test notification",
			args: args{
				title:   "Test Notification",
				message: "This is a test notification",
				link:    "https://joyfere.cn",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip the test if we're on an unsupported platform
			// Currently supporting macOS and Windows
			if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
				t.Skip("Skipping test on unsupported platform:", runtime.GOOS)
			}

			if err := sendNotification(tt.args.title, tt.args.message, tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("sendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
