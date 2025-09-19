package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Entry

func Init(serviceName string) {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.JSONFormatter{})
	baseLogger.SetOutput(os.Stdout)
	baseLogger.SetLevel(logrus.InfoLevel)

	// Dodaj servis ime kao default polje
	Log = baseLogger.WithField("service", serviceName)
}

func Info(msg string, fields ...logrus.Fields) {
	entry := Log
	if len(fields) > 0 {
		entry = Log.WithFields(fields[0])
	}
	entry.Info(msg)
}

func Error(msg string, err error, fields ...logrus.Fields) {
	entry := Log.WithError(err)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(msg)
}

func Debug(msg string, fields ...logrus.Fields) {
	entry := Log
	if len(fields) > 0 {
		entry = Log.WithFields(fields[0])
	}
	entry.Debug(msg)
}

func Warn(msg string, fields ...logrus.Fields) {
	entry := Log
	if len(fields) > 0 {
		entry = Log.WithFields(fields[0])
	}
	entry.Warn(msg)
}

func WithRequestID(requestID string) *logrus.Entry {
	return Log.WithField("request_id", requestID)
}

func WithUserID(userID string) *logrus.Entry {
	return Log.WithField("user_id", userID)
}

func WithBlogID(blogID string) *logrus.Entry {
	return Log.WithField("blog_id", blogID)
}
