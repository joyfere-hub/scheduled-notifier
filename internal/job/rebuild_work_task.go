package job

import (
	"encoding/json"
	"fmt"
	"github.com/joyfere-hub/scheduled-notifier/notifier"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joyfere-hub/scheduled-notifier/internal/conf"
)

const (
	endpoint               = "https://api.task.work.rebuildsoft.com"
	loginUrl               = endpoint + "/client/1/user/login"
	getNotificationListUrl = endpoint + "/client/1/notification/list"
	getOrgListUrl          = endpoint + "/client/1/org/org-list"

	uiEndpoint = "https://task.work.rebuildsoft.com/#"
	taskPage   = uiEndpoint + "/task"
)

var logger = log.New(os.Stderr, "RebuildWorkTaskClient", log.LstdFlags)

type RebuildWorkTaskClient struct {
	Config        *conf.JobConfig
	LastFetchTime time.Time
	token         string
}

type baseResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type loginResult struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type notificationListResult struct {
	Count            int                   `json:"count"`
	NotificationList *[]notificationDetail `json:"notification_list"`
}
type notificationDetail struct {
	Id       int                  `json:"id"`
	Content  *notificationContent `json:"notification_content"`
	FromUser *UserDetail          `json:"from_user"`
}

type notificationContent struct {
	Id   int         `json:"id"`
	Data *taskDetail `json:"data"`
}

type taskDetail struct {
	Message string `json:"message"`
}

type UserDetail struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type OrgListResult struct {
	OrgList *[]orgDetail `json:"org_list"`
}

type orgDetail struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	register("rebuild_work_task", newClient)
}

func newClient(jobConfig *conf.JobConfig) (JobClient, error) {
	var err error
	token := jobConfig.Token
	if len(jobConfig.Token) == 0 {
		token, err = getToken(jobConfig.Username, jobConfig.Password)
		if err != nil {
			return nil, err
		}
	}
	return &RebuildWorkTaskClient{
		Config: jobConfig,
		token:  token,
	}, nil
}

func (c *RebuildWorkTaskClient) FetchMessages() (*notifier.Messages, error) {
	orgList, err := c.GetOrgList()
	if err != nil {
		return nil, err
	}
	if orgList == nil {
		return &notifier.Messages{}, nil
	}
	var messages notifier.Messages
	for _, org := range *orgList {
		notificationList, err := c.getNotificationList(org)
		if err != nil {
			return nil, err
		}
		for _, notification := range *notificationList {
			message, err := newMessage(&org, &notification)
			if err != nil {
				return nil, err
			}
			messages = append(messages, message)
		}
	}
	c.LastFetchTime = time.Now()
	return &messages, nil
}

func (c *RebuildWorkTaskClient) getNotificationList(org orgDetail) (*[]notificationDetail, error) {
	resp, err := http.PostForm(getNotificationListUrl, map[string][]string{
		"token":       {c.token},
		"org_root_id": {strconv.Itoa(org.Id)},
		"page":        {"1"},
		"read":        {"0"},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get notification list fail! status code : %s", resp.Status)
	}
	result, err := parseResponse[notificationListResult](resp)
	if err != nil {
		return nil, err
	}
	return result.NotificationList, nil
}

func newMessage(org *orgDetail, notification *notificationDetail) (notifier.Message, error) {
	if notification.Content == nil {
		return notifier.Message{}, fmt.Errorf("notification content is nil")
	}
	if notification.Content.Data == nil {
		return notifier.Message{}, fmt.Errorf("notification content data is nil")
	}
	return notifier.Message{
		Title:   fmt.Sprintf("组织「%s」成员「%s」创建了任务", org.Name, notification.FromUser.Name),
		Message: notification.Content.Data.Message,
		Link:    taskPage,
	}, nil
}

func (c *RebuildWorkTaskClient) GetOrgList() (*[]orgDetail, error) {
	resp, err := http.PostForm(getOrgListUrl, map[string][]string{
		"token": {c.token},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get org list fail! status code : %s", resp.Status)
	}
	result, err := parseResponse[OrgListResult](resp)
	if err != nil {
		return nil, err
	}
	return result.OrgList, nil
}

func getToken(username, password string) (string, error) {
	resp, err := http.PostForm(loginUrl, map[string][]string{
		"phone":    {username},
		"password": {password},
		"token":    {},
	})
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get rebuild work task token fail! status code: %s", resp.Status)
	}
	result, err := parseResponse[loginResult](resp)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}

func parseResponse[T any](resp *http.Response) (*T, error) {
	if resp == nil {
		return nil, fmt.Errorf("http response is nil")
	}
	defer resp.Body.Close()

	var result baseResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("http response is fail, internal code : %d", result.Code)
	}

	return &result.Data, nil
}
