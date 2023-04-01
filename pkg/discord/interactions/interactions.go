package interactions

import (
  "fmt"

	"github.com/bwmarrin/discordgo"
)

type InteractionLogger struct {
}

func (l *InteractionLogger) SendDeferredInteractionMessage(s *discordgo.Session, i *discordgo.Interaction) error {
	if err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	}); err != nil {
    return fmt.Errorf("s.InteractionRespond: %v", err)
	}
  return nil
}

func (l *InteractionLogger) SendInteractionMessage(s *discordgo.Session, i *discordgo.Interaction, msg string) error {
	if err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	}); err != nil {
    return fmt.Errorf("s.InteractionRespond: %v", err)
	}
  return nil
}

func (l *InteractionLogger) SendInteractionFiles(s *discordgo.Session, i *discordgo.Interaction, files []*discordgo.File) error {
	if err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: files,
		},
	}); err != nil {
    return fmt.Errorf("s.InteractionRespond: %v", err)
	}
  return nil
}

func (l *InteractionLogger) SendEditedInteractionMessage(s *discordgo.Session, i *discordgo.Interaction, msg string) error {
	if _, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &msg,
	}); err != nil {
    return fmt.Errorf("s.InteractionRespond: %v", err)
	}
  return nil
}

func (l *InteractionLogger) SendEditedInteractionFiles(s *discordgo.Session, i *discordgo.Interaction, files []*discordgo.File) error {
	if _, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Files: files,
	}); err != nil {
    return fmt.Errorf("s.InteractionRespond: %v", err)
	}
  return nil
}

func FindCommandDataByName(name string, options []*discordgo.ApplicationCommandInteractionDataOption) *discordgo.ApplicationCommandInteractionDataOption {
	for _, d := range options {
		if d.Name == name {
			return d
		}
		if o := FindCommandDataByName(name, d.Options); o != nil {
			return o
		}
	}
	return nil
}
