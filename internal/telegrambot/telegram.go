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

	bot struct {
		cli    *tgbotapi.BotAPI
		chatID int64
	}
)

func NewBot(config Config) (*bot, error) {
	cli, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("new tg bot: %w", err)
	}

	return &bot{
		cli:    cli,
		chatID: config.ChatID,
	}, nil
}

func (b *bot) Send(text string) error {
	msg := tgbotapi.NewMessage(b.chatID, text)
	_, err := b.cli.Send(msg)
	return err
}
