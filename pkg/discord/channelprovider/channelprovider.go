package channelprovider

import (
	"time"

	"github.com/bwmarrin/discordgo"
	types "github.com/thompsonja/discord_bots_lib/pkg/types/channelprovider"
)

type ChannelProvider struct {
	Discord *discordgo.Session
}

func DiscordUserToUser(duser *discordgo.User) *types.User {
	return &types.User{
		Bot:      duser.Bot,
		ID:       duser.ID,
		Username: duser.Username,
	}
}

func (d *ChannelProvider) Channel(channelID string) (*types.Channel, error) {
	channel, err := d.Discord.Channel(channelID)
	if err != nil {
		return nil, err
	}

	return &types.Channel{
		LastMessageID: channel.LastMessageID,
	}, nil
}

func (d *ChannelProvider) ChannelMessages(channelID string, limit int, beforeID string, afterID string, aroundID string) ([]*types.ChannelMessage, error) {
	dmsgs, err := d.Discord.ChannelMessages(channelID, limit, beforeID, afterID, aroundID)
	if err != nil {
		return []*types.ChannelMessage{}, err
	}

	msgs := []*types.ChannelMessage{}
	for _, m := range dmsgs {
		mentions := []*types.User{}
		for _, mention := range m.Mentions {
			mentions = append(mentions, DiscordUserToUser(mention))
		}
		msgs = append(msgs, &types.ChannelMessage{
			Author:    DiscordUserToUser(m.Author),
			ChannelID: m.ChannelID,
			Content:   m.Content,
			ID:        m.ID,
			Mentions:  mentions,
		})
	}

	return msgs, nil
}

func (d *ChannelProvider) ChannelTyping(channelID string) error {
	return d.Discord.ChannelTyping(channelID)
}

func (d *ChannelProvider) SnowflakeTimestamp(ID string) (time.Time, error) {
	return discordgo.SnowflakeTimestamp(ID)
}
