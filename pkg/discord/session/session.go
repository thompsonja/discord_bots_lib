package session

import (
  "context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thompsonja/discord_bots_lib/pkg/gcp/secrets"
)

func GetSession(ctx context.Context, secretKey, projectID string) (*discordgo.Session, error) {
	botKey, err := secrets.GetLatestSecretValue(ctx, secretKey, projectID)
	if err != nil {
		return nil, fmt.Errorf("error getting bot key: %v", err)
	}

	d, err := discordgo.New("Bot " + botKey)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}
	return d, nil
}
