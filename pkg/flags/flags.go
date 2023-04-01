package flags

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func GetBotKey(botKeyFlag, botKeyFileFlag string) (string, error) {
	botKeyEnvVar := os.Getenv("BOT_TOKEN")
	botKeyFileEnvVar := os.Getenv("BOT_TOKEN_FILE")

	if (botKeyFlag != "") && (botKeyEnvVar != "") && (botKeyFlag != botKeyEnvVar) {
		return "", errors.New("--bot_token provided and BOT_TOKEN defined, and they differ")
	}

	if botKeyFlag != "" {
		return botKeyFlag, nil
	} else if botKeyEnvVar != "" {
		return botKeyEnvVar, nil
	}

	if (botKeyFileFlag != "") && (botKeyFileEnvVar != "") && (botKeyFileFlag != botKeyFileEnvVar) {
		return "", fmt.Errorf("--bot_token_file provided and BOT_TOKEN_FILE defined, and they differ (%s vs %s)", botKeyFileFlag, botKeyFileEnvVar)
	}

	botKeyFile := ""
	if botKeyFileFlag != "" {
		botKeyFile = botKeyFileFlag
	} else if botKeyFileEnvVar != "" {
		botKeyFile = botKeyFileEnvVar
	} else {
		botKeyFile = "/etc/discord/bot_secret"
	}

	info, err := os.Stat(botKeyFile)
	if os.IsNotExist(err) || info.IsDir() {
		return "", fmt.Errorf("-bot_token not provided, and bot_token_file %s could not be found", botKeyFile)
	}

	buf, err := os.ReadFile(botKeyFile)
	if err != nil {
		return "", fmt.Errorf("error reading bot token file: %v", err)
	}
	return strings.TrimSpace(string(buf)), nil
}
