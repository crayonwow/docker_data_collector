package config

import (
	"fmt"
	"os"
	"strconv"

	"docker_data_collector/internal/sender/telegrambot"
)

type (
	config struct {
		TG telegrambot.Config
	}
)

func newConfig() (*config, error) {
	tg, err := tgConfig()
	return &config{tg}, err
}

func tgConfig() (telegrambot.Config, error) {
	chatIDRaw := os.Getenv("TG_CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDRaw, 10, 64)
	if err != nil {
		return telegrambot.Config{}, fmt.Errorf("parese int %s: %w", chatIDRaw, err)
	}
	return telegrambot.Config{
		Token:  os.Getenv("TG_BOT_KEY"),
		ChatID: chatID,
	}, err
}

func TGConfig(c *config) telegrambot.Config {
	return c.TG
}
