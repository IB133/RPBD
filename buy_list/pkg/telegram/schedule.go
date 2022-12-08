package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/IB133/RPBD/buy_list/pkg/config"
	"github.com/IB133/RPBD/buy_list/pkg/db"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) StartScheduler(cnf *config.Config) {
	cnf.BuyListSched = gocron.NewScheduler(time.Local)
	cnf.FridgeSched = gocron.NewScheduler(time.Local)
	users := db.UsersList(*b.que)
	for _, u := range users {
		id := u.Id
		chatId := u.Chat_id

		buyStart, err := cnf.BuyListSched.Every(1).Seconds().Do(createBuyListScheduler, b, chatId, id)
		if err != nil {
			log.Fatal(err)
		}
		buyStart.LimitRunsTo(1)

		cnf.BuyListSched.Every(1).Day().At("08:00").Do(createBuyListScheduler, b, chatId, id)
		cnf.BuyListSched.StartAsync()

		fridgeStart, err := cnf.FridgeSched.Every(1).Seconds().Do(createFridgeScheduler, b, chatId, id)
		if err != nil {
			log.Fatal(err)
		}
		fridgeStart.LimitRunsTo(1)

		cnf.FridgeSched.Every(1).Day().At("08:00").Do(createFridgeScheduler, b, chatId, id)
		cnf.FridgeSched.StartAsync()
	}
}

func createBuyListScheduler(b *Bot, chatId int64, userId int) {
	list := db.SchedulerBuyList(userId, *b.que)
	if list == nil {
		return
	}
	str := "Список покупок на сегодня:\n"
	for _, v := range list {
		str += fmt.Sprintf("%s %f\n", v.Prod_name, v.Weight)
	}
	msg := tgbotapi.NewMessage(chatId, str)
	if _, err := b.bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func createFridgeScheduler(b *Bot, chatId int64, userId int) {
	var str string
	list := db.SchedulerFridge(userId, *b.que)
	if list == nil {
		return
	}
	str = "Сегодня выходит срок годности на данные продукты:\n"
	for _, v := range list {
		switch v.Status {
		case "stored":
			v.Status = "хранится"
		case "opened":
			v.Status = "открыт"
		}
		str += fmt.Sprintf("%s %s %s\n", v.Prod_name, v.Status, v.Experitation_date.Format("2006-01-02"))
	}
	msg := tgbotapi.NewMessage(chatId, str)
	if _, err := b.bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}
