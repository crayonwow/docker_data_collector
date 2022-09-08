package config

import (
	"fmt"
	"os"
	"strconv"

	"docker_data_collector/internal/telegrambot"
)

type (
	Config struct {
		TG telegrambot.Config
	}
)

func New() (Config, error) {
	tg, err := tgConfig()
	return Config{tg}, err
}

func tgConfig() (telegrambot.Config, error) {
	chatID, err := strconv.ParseInt(os.Getenv("TG_CHAT_ID"), 10, 64)
	if err != nil {
		return telegrambot.Config{}, fmt.Errorf("parese int %s: %w", chatID, err)
	}
	return telegrambot.Config{
		Token:  os.Getenv("TG_BOT_KEY"),
		ChatID: chatID,
	}, err
}
