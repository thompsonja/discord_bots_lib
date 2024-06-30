package channellogger

import (
	"io"
	"log"
)

type TestChannelLogger struct{}

func (t *TestChannelLogger) SendMessage(channelID, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (t *TestChannelLogger) SendFile(channelID, name string, reader io.Reader) {
	log.Printf("Sent file %s\n", name)
}
