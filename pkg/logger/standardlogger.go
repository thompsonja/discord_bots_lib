package logger

import (
  "log"
)

type standardlogger struct{}

func (l *standardlogger) Debug(args ...interface{}) {
	log.Println(args...)
}

func (l *standardlogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *standardlogger) Info(args ...interface{}) {
	log.Println(args...)
}

func (l *standardlogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *standardlogger) Warning(args ...interface{}) {
	log.Println(args...)
}

func (l *standardlogger) Warningf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *standardlogger) Error(args ...interface{}) {
	log.Println(args...)
}

func (l *standardlogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *standardlogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (l *standardlogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
