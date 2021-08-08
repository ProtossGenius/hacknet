package hnlog

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// NewLogger new logger.
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter.(*logrus.JSONFormatter).DisableTimestamp = false // remove timestamp from test output.
	logger.Formatter.(*logrus.JSONFormatter).TimestampFormat = "2006-01-02 15:04:05"
	logger.Out = os.Stdout

	return logger
}

var globalLog = NewLogger()

// DefaultLogger .
func DefaultLogger() *logrus.Logger {
	return globalLog
}

func logTime() *logrus.Entry {
	return globalLog.WithTime(time.Now())
}

// Info write log info.
func Info(info string, fields logrus.Fields) {
	logTime().WithFields(fields).Infoln(info)
}

// Warn write warning log.
func Warn(info string, fields logrus.Fields) {
	logTime().WithFields(fields).Warnln(info)
}

// Error write error log.
func Error(info string, fields logrus.Fields) {
	logTime().WithFields(fields).Errorln()
}
