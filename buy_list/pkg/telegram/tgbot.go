package telegram

import (
	"fmt"
	"log"
	"os"

	"github.com/IB133/RPBD/buy_list/pkg/config"
	"github.com/IB133/RPBD/buy_list/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	cnf *config.Config
	que *db.DB
}

func NewBot(cfg *config.Config) *Bot {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Env file doesnt exist: %s", err)
	}
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := db.NewConnect(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DBNAME"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBUSER")))
	if err != nil {
		log.Fatal()
	}
	cfg.Keyboards = *config.NewKeyboard()
	return &Bot{
		bot: botAPI,
		cnf: cfg,
		que: conn,
	}
}

func (b *Bot) Start() error {
	b.bot.Debug = true
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.StartScheduler(b.cnf)

	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.Message.IsCommand() {
			if err := b.handleCommands(update.Message, &msg); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if b.cnf.UserInsert {
			if err := b.handleMessages(update.Message, &msg); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if err := b.handleKeyboard(update.Message, &msg); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
