package channellogger

import (
	"fmt"
	"io"
	"log"

	"github.com/bwmarrin/discordgo"
)

type ChannelLogger struct {
	discord *discordgo.Session
}

func New(discord *discordgo.Session) *ChannelLogger {
	return &ChannelLogger{
		discord: discord,
	}
}

func (l *ChannelLogger) SendMessage(channelID string, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if message == "" {
		return
	}
	if l.discord == nil {
		log.Println("Nil discord session when attempting to SendMessage")
	}
	_, err := l.discord.ChannelMessageSend(channelID, message)
	log.Println("Sending message to channel", channelID, ":", message)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

func (l *ChannelLogger) SendFile(channelID string, name string, reader io.Reader) {
	if l.discord == nil {
		log.Println("Nil discord session when attempting to SendMessage")
	}
	_, err := l.discord.ChannelFileSend(channelID, name, reader)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
