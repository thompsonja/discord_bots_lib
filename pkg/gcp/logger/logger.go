package logger

import (
	"fmt"
	"os"

	"cloud.google.com/go/logging"
)

type CloudLogger struct {
	client *logging.Client
	logger *logging.Logger
}

func New(client *logging.Client, name string) *CloudLogger {
	return &CloudLogger{
		client: client,
		logger: client.Logger(name),
	}
}

func (l *CloudLogger) log(payload string, severity logging.Severity) {
	l.logger.Log(logging.Entry{
		Payload:  payload,
		Severity: severity,
	})
}

func (l *CloudLogger) Debug(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Debug)
}

func (l *CloudLogger) Debugf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Debug)
}

func (l *CloudLogger) Info(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Info)
}

func (l *CloudLogger) Infof(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Info)
}

func (l *CloudLogger) Warning(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Warning)
}

func (l *CloudLogger) Warningf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Warning)
}

func (l *CloudLogger) Error(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Error)
}

func (l *CloudLogger) Errorf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Error)
}

func (l *CloudLogger) Fatal(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Critical)
	os.Exit(1)
}

func (l *CloudLogger) Fatalf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Critical)
	os.Exit(1)
}
