package channellogger

import (
	"log"
)

type TestChannelLogger struct{}

func (t *TestChannelLogger) SendMessage(channelID string, format string, v ...interface{}) {
	log.Printf(format, v...)
}
