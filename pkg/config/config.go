package config

import (
	"fmt"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
)

// Config is a type alias for the configuration object
type Config struct {
	*viper.Viper
}

var configFile string
var config *viper.Viper

func init() {
	c := viper.New()

	c.SetEnvPrefix("Q")
	c.AutomaticEnv()

	flag.StringVar(&configFile, "config", "", "")
	flag.Parse()
	c.BindPFlag("config", flag.Lookup("config"))

	if c.GetString("config") == "" {
		// Read from "default" configuration path
		c.SetConfigName("config")
		c.AddConfigPath("/etc/qurid")
		c.AddConfigPath("$HOME/.quarid")
		c.AddConfigPath(".")
	} else {
		c.SetConfigFile(c.GetString("config"))
	}

	if err := c.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Unable to read any configuration file: %s\n", err))
	}

	location, err := time.LoadLocation(c.GetString("timezone"))
	if err == nil {
		c.Set("timezone", location)
	} else {
		c.Set("timezone", time.UTC)
	}

	config = c
}

// Get returns the global configuration object
func Get() *viper.Viper {
	return config
}
