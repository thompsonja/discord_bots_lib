package session

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thompsonja/discord_bots_lib/pkg/gcp/secrets"
)

func GetSession(secretKey, projectID string) (*discordgo.Session, error) {
	botKey, err := secrets.GetLatestSecretValue(secretKey, projectID)
	if err != nil {
		return nil, fmt.Errorf("error getting bot key: %v", err)
	}

	d, err := discordgo.New("Bot " + botKey)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}
	return d, nil
}
