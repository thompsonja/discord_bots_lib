package channellogger

import "io"

type ChannelLogger interface {
	SendMessage(channelID string, format string, v ...interface{})
	SendFile(channelID string, name string, reader io.Reader)
}
