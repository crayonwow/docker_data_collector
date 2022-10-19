package telegrambot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	Config struct {
		Token  string
		ChatID int64
	}

	Bot struct {
		cli    *tgbotapi.BotAPI
		chatID int64
	}
)

func newBot(config Config) (*Bot, error) {
	cli, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("new tg bot: %w", err)
	}

	return &Bot{
		cli:    cli,
		chatID: config.ChatID,
	}, nil
}

func (b *Bot) Send(text string) error {
	msg := tgbotapi.NewMessage(b.chatID, text)
	_, err := b.cli.Send(msg)
	return err
}
