package channellogger

import (
	"fmt"
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
	_, err := l.discord.ChannelMessageSend(channelID, message)
	log.Println("Sending message to channel", channelID, ":", message)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
