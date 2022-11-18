package tgbot

import (
	"log"

	"github.com/IB133/RPBD/buy_list/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	messages  config.Responses
	keyboards config.Keyboards
}

func NewBot(b *tgbotapi.BotAPI, cfg *config.Config) *Bot {
	k := cfg.NewKeyboard()
	return &Bot{
		bot:       b,
		messages:  cfg.Responses,
		keyboards: k,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.IsCommand() {
			if err := b.handleCommands(update.Message); err != nil {
				log.Fatal()
			}
			continue
		}
		if err := b.handleIncorrect(update.Message); err != nil {
			log.Fatal()
		}
		if err := b.handleMessages(update.Message); err != nil {
			log.Fatal()
		}
	}
	return nil
}
