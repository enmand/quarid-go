package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"
)

type Config struct {
	Nick         string         `json:"nick"`
	User         string         `json:"user"`
	TimezoneName string         `json:"timezone"`
	Timezone     *time.Location `json:"-"`

	Server string `json:"server"`
	TLS    struct {
		Verify bool `json:"verify"`
		Enable bool `json:"enable"`
	} `json:"tls"`
	Channels []string `json:"channels"`
	Admins   []string `json:"admins"`

	Log logType `json:"log"`

	PluginBlacklist []string `json:"plugin_blacklist"`

	Logger *logrus.Logger `json:"-"`
}

func NewConfig(configFilePath string) (*Config, error) {
	var configData []byte
	var config Config
	var err error

	// Read configuration file
	configData, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("Could not load configuration: %s", err)
	}

	// Load configuration struct
	if err = json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("Could not parse JSON: %s", err)
	}

	w, err := logWriter(config.Log)
	if err != nil {
		panic(fmt.Sprintf("Cannot configure log writer: %s", err))
	}
	config.Logger = NewLogger(w, config.Log)

	config.Logger.Printf("Loading timezone...")
	// Load timezone
	location, err := time.LoadLocation(config.TimezoneName)
	if err == nil {
		config.Timezone = location
	} else {
		config.Timezone = time.UTC
	}

	config.Logger.Printf("Using timezone: '%s'\n", config.Timezone.String())

	return &config, nil
}
