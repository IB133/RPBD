package main

import (
	"log"

	"github.com/IB133/RPBD/buy_list/pkg/config"
	tgbot "github.com/IB133/RPBD/buy_list/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	botAPI, err := tgbotapi.NewBotAPI("5434976758:AAEoFg0TIO7aCCvoL2s3MK75NmGgfB0lbgI")
	if err != nil {
		log.Fatal(err)
	}
	botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)
	bot := tgbot.NewBot(botAPI, cfg)
	bot.Start()
}
