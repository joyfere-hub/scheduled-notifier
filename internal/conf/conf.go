package conf

import (
	"encoding/json"
	"github.com/segmentfault/pacman/contrib/conf/viper"
	"os"
	"path/filepath"
)

const (
	DefConfigFilePath    = "/.local/scheduled-notifier.yaml"
	DefConfigFilePathEnv = "SCHEDULED_NOTIFIER_CONFIG_FILE_PATH"
)

type BaseConfig struct {
	Jobs *[]JobConfig `json:"jobs" yaml:"jobs"`
}

func (c *BaseConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

type JobConfig struct {
	Name     string `json:"name" yaml:"name"`
	Interval string `json:"interval" yaml:"interval"`

	// base auth
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	// token auth
	Token string `json:"token" yaml:"token"`
}

func (c *JobConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func ReadConfig(configFilePath string) (c *BaseConfig, err error) {
	if len(configFilePath) == 0 {
		return readDefConfig()
	}
	return readConfig(configFilePath)
}

func readConfig(configFilePath string) (c *BaseConfig, err error) {
	c = &BaseConfig{}
	config, err := viper.NewWithPath(configFilePath)
	if err != nil {
		return nil, err
	}
	if err = config.Parse(&c); err != nil {
		return nil, err
	}
	return c, nil
}

func readDefConfig() (c *BaseConfig, err error) {
	envConfigFilePath := os.Getenv(DefConfigFilePathEnv)
	if len(envConfigFilePath) > 0 {
		return readConfig(envConfigFilePath)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return readConfig(filepath.Join(home, DefConfigFilePath))
}
