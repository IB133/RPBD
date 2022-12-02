package main

import (
	"log"

	"github.com/IB133/RPBD/buy_list/pkg/config"
	tgbot "github.com/IB133/RPBD/buy_list/pkg/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot := tgbot.NewBot(cfg)
	bot.Start()
}
