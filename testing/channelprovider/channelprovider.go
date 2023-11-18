package channelprovider

import (
	"time"

	types "github.com/thompsonja/discord_bots_lib/pkg/types/channelprovider"
)

type TestChannelProvider struct{}

func (t *TestChannelProvider) Channel(channelID string) (*types.Channel, error) {
	return &types.Channel{
		LastMessageID: channelID,
	}, nil
}

func (t *TestChannelProvider) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*types.ChannelMessage, error) {
	return []*types.ChannelMessage{
		{
			Author: &types.User{
				Bot: false,
				ID:  channelID,
			},
			ChannelID: channelID,
			Content:   "Test content",
			ID:        beforeID,
		},
	}, nil
}

func (t *TestChannelProvider) ChannelTyping(channelID string) error {
	return nil
}

func (t *TestChannelProvider) SnowflakeTimestamp(ID string) (time.Time, error) {
	return time.Now(), nil
}
