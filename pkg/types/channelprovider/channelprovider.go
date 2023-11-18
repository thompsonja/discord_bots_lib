package channelprovider

import "time"

type Channel struct {
	LastMessageID string
}

type User struct {
	Bot      bool
	ID       string
	Username string
}

type ChannelMessage struct {
	Author    *User
	ChannelID string
	Content   string
	ID        string
	Mentions  []*User
}

type ChannelProvider interface {
	Channel(channelID string) (*Channel, error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*ChannelMessage, error)
	ChannelTyping(channelID string) error
	SnowflakeTimestamp(ID string) (time.Time, error)
}
