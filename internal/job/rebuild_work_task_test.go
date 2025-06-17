package job

import (
	"fmt"
	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
	"github.com/joyfere-hub/scheduled-notifier/notifier"
	"testing"
)

func Test_newMessage(t *testing.T) {
	tests := []struct {
		name          string
		org           *orgDetail
		notification  *notificationDetail
		wantErr       bool
		errorContains string
	}{
		{
			name: "valid message",
			org: &orgDetail{
				Id:   1,
				Name: "Test Org",
			},
			notification: &notificationDetail{
				Id: 1,
				Content: &notificationContent{
					Id: 1,
					Data: &taskDetail{
						Message: "Test Message",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nil content",
			org: &orgDetail{
				Id:   1,
				Name: "Test Org",
			},
			notification: &notificationDetail{
				Id:      1,
				Content: nil,
			},
			wantErr:       true,
			errorContains: "notification content is nil",
		},
		{
			name: "nil data",
			org: &orgDetail{
				Id:   1,
				Name: "Test Org",
			},
			notification: &notificationDetail{
				Id: 1,
				Content: &notificationContent{
					Id:   1,
					Data: nil,
				},
			},
			wantErr:       true,
			errorContains: "notification content data is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newMessage(tt.org, tt.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("newMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errorContains != "" {
				if err.Error() != tt.errorContains {
					t.Errorf("newMessage() error = %v, want error containing %v", err, tt.errorContains)
				}
			}
		})
	}
}

func Test_getToken(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				username: "phone",
				password: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getToken(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("getToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("token: %s\n", got)
		})
	}
}

func TestRebuildWorkTaskClient_FetchMessages(t *testing.T) {
	tests := []struct {
		name    string
		conf    *conf.JobConfig
		want    *[]notifier.Message
		wantErr bool
	}{
		{
			name: "base",
			conf: &conf.JobConfig{
				Token: "token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.conf)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
				return
			}
			got, err := c.FetchMessages()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(*got) == 0 {
				t.Log("no message")
				return
			}
			for i, message := range *got {
				t.Logf("message %d: %s", i, message.String())
				err := message.Send()
				if err != nil {
					return
				}
			}
		})
	}
}
