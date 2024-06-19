package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init(appEnv string) *logrus.Logger {
	var logger logrus.Logger
	switch appEnv {
	case "local":
		logger = logrus.Logger{
			Level: logrus.DebugLevel,
			Out:   os.Stderr,
			Formatter: &logrus.TextFormatter{
				DisableLevelTruncation: true,
				TimestampFormat:        "02.01.2006 Mon 15:04:05",
			},
			ReportCaller: false,
		}
	case "prod":
		logger = logrus.Logger{
			Level: logrus.InfoLevel,
			Out:   os.Stderr,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "02.01.2006 Mon 15:04:05",
			},
			ReportCaller: true,
		}
	}
	return &logger
}
