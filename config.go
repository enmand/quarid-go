package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Nick         string         `json:"nick"`
	User         string         `json:"user"`
	TimezoneName string         `json:"timezone"`
	Timeozone    *time.Location `json:"-"`

	Server string `json:"server"`
	TLS    struct {
		Verify bool `json:"verify"`
		Enable bool `json:"enable"`
	} `json:"tls"`
	Channels []string `json:"channels"`
	Admins   []string `json:"admins"`

	Plugins []string `json:"plugins"`

	LogFile interface{} `json:"log_file"`
	Log     *log.Logger `json:"-"`
}

func loadConfig(configFilePath string) (*Config, error) {
	var configData []byte
	var config Config
	var err error
	config.Log = setupLogging(true)

	// Read configuration file
	config.Log.Printf("Reading configuration file...\n")
	configData, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("Could not load configuration: %s", err)
	}

	// Load configuration struct
	config.Log.Printf("Unmarchalling JSON configuration file...\n")
	if err = json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("Could not parse JSON: %s", err)
	}
	log.Printf("Setting up \\'real\\' logging...\n")
	config.Log = setupLogging(config.LogFile)

	config.Log.Printf("Loading timezone...")
	// Load timezone
	location, err := time.LoadLocation(config.TimezoneName)
	if err == nil {
		config.Timeozone = location
	} else {
		config.Timeozone = time.UTC
	}

	config.Log.Printf("Using timezone: '%s'\n", config.Timeozone.String())

	return &config, nil
}

func setupLogging(loggingFile interface{}) *log.Logger {
	// Configure logger
	logWritter := os.Stderr

	if logFilePath, ok := loggingFile.(string); ok {
		logFile, err := os.Open(logFilePath)
		if err != nil {
			fmt.Printf("Could not load log file\n")
			defer logFile.Close()
		} else {
			logWritter = logFile
		}
	}

	return log.New(logWritter, "", log.LstdFlags)

}
