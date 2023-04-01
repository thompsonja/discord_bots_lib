package logger

import (
  "fmt"
	"os"

	"cloud.google.com/go/logging"
)

type cloudlogger struct {
	client *logging.Client
	logger *logging.Logger
}

func (l *cloudlogger) log(payload string, severity logging.Severity) {
	l.logger.Log(logging.Entry{
		Payload: payload,
		Severity: severity,
	})
}

func (l *cloudlogger) Debug(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Debug)
}

func (l *cloudlogger) Debugf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Debug)
}

func (l *cloudlogger) Info(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Info)
}

func (l *cloudlogger) Infof(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Info)
}

func (l *cloudlogger) Warning(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Warning)
}

func (l *cloudlogger) Warningf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Warning)
}

func (l *cloudlogger) Error(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Error)
}

func (l *cloudlogger) Errorf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Error)
}

func (l *cloudlogger) Fatal(args ...interface{}) {
	l.log(fmt.Sprint(args...), logging.Critical)
	os.Exit(1)
}

func (l *cloudlogger) Fatalf(format string, args ...interface{}) {
	l.log(fmt.Sprintf(format, args...), logging.Critical)
	os.Exit(1)
}
