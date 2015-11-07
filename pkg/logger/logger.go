package logger

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/enmand/quarid-go/pkg/config"
)

// Log can be logged to using Sirusen Logrus
var Log *log.Logger

func init() {
	c := config.Get()
	l := c.GetInt("log.level")

	Log = log.New()
	Log.Out = os.Stderr
	Log.Level = log.Level(l)
}
