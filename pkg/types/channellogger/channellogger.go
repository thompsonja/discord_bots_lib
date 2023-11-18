package channellogger

type ChannelLogger interface {
	SendMessage(channelID string, format string, v ...interface{})
}
