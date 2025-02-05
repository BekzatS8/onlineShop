package logging

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Настройка формата логов (JSON)
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Настройка вывода в файл
	logFile, err := os.OpenFile("server_logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(logFile)
	} else {
		Log.SetOutput(os.Stdout)
		Log.Warn("Failed to log to file, using default stdout")
	}

	Log.SetLevel(logrus.InfoLevel)
}
