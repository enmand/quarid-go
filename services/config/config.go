package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/enmand/quarid-go/services/logger"
)

type Config struct {
	Nick         string         `json:"nick"`
	Ident        string         `json:"user"`
	TimezoneName string         `json:"timezone"`
	Timezone     *time.Location `json:"-"`

	Server string `json:"server"`
	TLS    struct {
		Verify bool `json:"verify"`
		Enable bool `json:"enable"`
	} `json:"tls"`
	Channels []string `json:"channels"`
	Admins   []string `json:"admins"`

	Log logger.LogType `json:"log"`

	PluginsDirs []string `json:"plugins_dirs"`

	Logger *logrus.Logger `json:"-"`
}

func New(configFilePath string) (*Config, error) {
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

	w, err := logger.LogWriter(config.Log)
	if err != nil {
		panic(fmt.Sprintf("Cannot configure log writer: %s", err))
	}
	config.Logger = logger.New(w, config.Log)

	config.Logger.Printf("Loading timezone...")
	// Load timezone
	location, err := time.LoadLocation(config.TimezoneName)
	if err == nil {
		config.Timezone = location
	} else {
		config.Timezone = time.UTC
	}

	config.Logger.Printf("Using timezone: '%s'", config.Timezone.String())

	return &config, nil
}
