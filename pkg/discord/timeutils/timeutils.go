package timeutils

import (
  "time"

  "github.com/bwmarrin/discordgo"
  "github.com/pkg/errors"
)

type TimestampProvider interface {
  SnowflakeTimestamp(ID string) (time.Time, error)
}

func GetMessageLocalTimestamp(discord TimestampProvider, messageID string) (time.Time, error) {
  t, err := discordgo.SnowflakeTimestamp(messageID)
  if err != nil {
    return time.Time{}, errors.Wrap(err, "discord.SnowflakeTimestamp")
  }
  return t.Local(), nil
}
