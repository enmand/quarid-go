package logger

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

type LogType struct {
	File  string `json:"file"`
	Level uint   `json:"level"`
}

func New(f io.Writer, l LogType) *log.Logger {
	// Configure logger
	log.SetOutput(f)
	log.SetLevel(log.Level(l.Level))
	return log.New() //log.New(logWritter, "", log.LstdFlags)
}

func LogWriter(l LogType) (io.Writer, error) {
	switch l.File {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return os.Open(l.File)
	}
}
