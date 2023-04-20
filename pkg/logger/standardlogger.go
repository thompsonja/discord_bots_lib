package logger

import (
	"log"
)

type StandardLogger struct{}

func (l *StandardLogger) Debug(args ...interface{}) {
	log.Println(args...)
}

func (l *StandardLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *StandardLogger) Info(args ...interface{}) {
	log.Println(args...)
}

func (l *StandardLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *StandardLogger) Warning(args ...interface{}) {
	log.Println(args...)
}

func (l *StandardLogger) Warningf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *StandardLogger) Error(args ...interface{}) {
	log.Println(args...)
}

func (l *StandardLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *StandardLogger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (l *StandardLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
