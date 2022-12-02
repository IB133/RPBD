package tgbot

import (
	"fmt"
	"strings"

	"github.com/IB133/RPBD/buy_list/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(mes *tgbotapi.Message, msg *tgbotapi.MessageConfig) error {
	switch mes.Command() {
	case "start":
		if _, err := b.que.GetUserByUsername(mes.From.UserName); err == nil {
			msg.Text = b.cnf.Start
			msg.ReplyMarkup = b.cnf.Main
			break
		}
		b.que.AddUser(mes.From.UserName, int(mes.Chat.ID))
		msg.Text = b.cnf.Start
		msg.ReplyMarkup = b.cnf.Main
	default:
		msg.Text = "Неизвестная команда"
	}
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleKeyboard(mes *tgbotapi.Message, msg *tgbotapi.MessageConfig) error {
	switch mes.Text {
	case b.cnf.Main.Keyboard[0][0].Text:
		msg.Text = "Добавление в список покупок"
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.AddToBuyList = true
		b.cnf.UserInsert = true

	// case b.keyboards.Cancel.Keyboard[0][0].Text:
	// 	msg.ReplyMarkup = b.keyboards.Main

	case b.cnf.Main.Keyboard[0][1].Text:
		msg.Text = "Выберите откуда добавить продукт"
		msg.ReplyMarkup = b.cnf.Keyboards.BuyOrNew
		b.cnf.Current = b.cnf.BuyOrNew

	case b.cnf.Main.Keyboard[1][0].Text:
		msg.Text = db.StoredProductList(mes.From.UserName, *b.que, *b.cnf)
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.OpenProduct = true
		b.cnf.UserInsert = true

	case b.cnf.Main.Keyboard[1][1].Text:
		msg.Text = fmt.Sprintf("Выберите продукт из списка\n%s", db.FridgeList(mes.From.UserName, *b.que, *b.cnf))
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.ChangeStatus = true
		b.cnf.UserInsert = true

	case b.cnf.Main.Keyboard[1][2].Text:
		msg.Text = "Введите две даты, между которыми нужно получить статистику."
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.GetStatistic = true
		b.cnf.UserInsert = true

	case b.cnf.Main.Keyboard[2][0].Text:
		msg.Text = db.FridgeList(mes.From.UserName, *b.que, *b.cnf)

	case b.cnf.Main.Keyboard[2][1].Text:
		msg.Text = db.UsedProcutList(mes.From.UserName, *b.que, *b.cnf)

	case b.cnf.BuyOrNew.Keyboard[0][0].Text:
		b.cnf.AddToFridgeFromBuyList = true
		b.cnf.UserInsert = true
		str, err := db.GetBuyList(mes.From.UserName, *b.que)
		if err != nil {
			msg.Text = err.Error()
			break
		}
		msg.Text = "Выберите один из продуктов\n" + str

	case b.cnf.BuyOrNew.Keyboard[0][1].Text:
		b.cnf.AddToFridge = true
		b.cnf.UserInsert = true
		msg.Text = "Введите продукт"

	case b.cnf.BuyOrNew.Keyboard[0][2].Text:
		msg.ReplyMarkup = b.cnf.Main
		b.cnf.Current = b.cnf.Main

	default:
		msg.Text = "Выберите пункт из меню"
		msg.ReplyMarkup = b.cnf.Current
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleMessages(mes *tgbotapi.Message, msg *tgbotapi.MessageConfig) error {
	switch {
	case mes.Text == "Отмена":
		b.cnf.AddToFridge = false
		b.cnf.AddToBuyList = false
		b.cnf.UserInsert = false
		msg.ReplyMarkup = b.cnf.Main
		b.cnf.Current = b.cnf.Main

	case b.cnf.AddToBuyList:
		str := strings.Split(mes.Text, " ")
		if len(str) != 3 {
			msg.Text = "Неверный ввод"
			break
		}
		u, err := b.que.GetUserByUsername(mes.From.UserName)
		if err != nil {
			msg.Text = err.Error()
			break
		}
		if err := b.que.AddProductToBuyList(u.Id, str[0], str[1], str[2]); err != nil {
			msg.Text = err.Error()
			break
		}
		msg.Text = "Заебсиь чотко"
		b.cnf.AddToBuyList = false
		b.cnf.UserInsert = false

	case b.cnf.AddToFridge:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = "Неверный ввод"
			break
		}
		if err := db.AddProductToFridge(str[0], mes.From.UserName, str[1], *b.que); err != nil {
			msg.Text = err.Error()
			break
		}
		msg.Text = "Заебсиь чотко"
		msg.ReplyMarkup = b.cnf.BuyOrNew
		b.cnf.AddToFridge = false
		b.cnf.UserInsert = false

	case b.cnf.AddToFridgeFromBuyList:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = "Неверный ввод"
			break
		}
		if err := db.AddProductToFridgeFromBuyList(str[0], mes.From.UserName, str[1], *b.que); err != nil {
			msg.Text = err.Error()
			break
		}
		msg.Text = "Заебсиь чотко"
		msg.ReplyMarkup = b.cnf.BuyOrNew
		b.cnf.AddToFridgeFromBuyList = false
		b.cnf.UserInsert = false

	case b.cnf.OpenProduct:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = b.cnf.ErrorInsert
			break
		}
		msg.Text = db.OpenProduct(mes.From.UserName, str[0], str[1], *b.que, *b.cnf)
		b.cnf.OpenProduct = false
		b.cnf.UserInsert = false

	case b.cnf.ChangeStatus:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = b.cnf.ErrorInsert
			break
		}
		msg.Text = db.ChangeStatus(mes.From.UserName, str[0], str[1], *b.que, *b.cnf)
		b.cnf.ChangeStatus = false
		b.cnf.UserInsert = false

	case b.cnf.GetStatistic:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = b.cnf.ErrorInsert
			break
		}
		msg.Text = db.GetStats(mes.From.UserName, str[0], str[1], *b.que, *b.cnf)
		b.cnf.GetStatistic = false
		b.cnf.UserInsert = false

	}
	_, err := b.bot.Send(msg)
	return err
}
