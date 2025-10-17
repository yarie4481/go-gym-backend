package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()

	// Use JSON formatter
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)

	// Ensure logs directory exists
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
	}

	// Output to file
	file, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Warn("Failed to log to file, using default stdout")
	}

	// Also log to stdout
	Log.SetOutput(os.Stdout)
}
