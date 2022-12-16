package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(mes *tgbotapi.Message, msg *tgbotapi.MessageConfig) error {
	switch mes.Command() {
	case "start":
		msg.Text = b.que.AddUser(mes.From.UserName, mes.Chat.ID, *b.cnf)
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
		msg.Text = b.cnf.AddingBuyList
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.AddToBuyList = true
		b.cnf.UserInsert = true

	// case b.keyboards.Cancel.Keyboard[0][0].Text:
	// 	msg.ReplyMarkup = b.keyboards.Main

	case b.cnf.Main.Keyboard[0][1].Text:
		msg.Text = b.cnf.AddingFridgeList
		msg.ReplyMarkup = b.cnf.Keyboards.BuyOrNew
		b.cnf.Current = b.cnf.BuyOrNew

	case b.cnf.Main.Keyboard[1][0].Text:
		msg.Text = b.que.StoredProductList(mes.From.UserName, *b.cnf)
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.OpenProduct = true
		b.cnf.UserInsert = true

	case b.cnf.Main.Keyboard[1][1].Text:
		msg.Text = fmt.Sprintf("%s\n%s", b.cnf.StatusChange, b.que.FridgeList(mes.From.UserName, *b.cnf))
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.ChangeStatus = true
		b.cnf.UserInsert = true

	case b.cnf.Main.Keyboard[1][2].Text:
		msg.Text = b.cnf.GetStats
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.GetStatistic = true
		b.cnf.UserInsert = true

	case b.cnf.Main.Keyboard[2][0].Text:
		msg.Text = b.que.FridgeList(mes.From.UserName, *b.cnf)

	case b.cnf.Main.Keyboard[2][1].Text:
		msg.Text = b.que.UsedProcutList(mes.From.UserName, *b.cnf)

	case b.cnf.BuyOrNew.Keyboard[0][0].Text:
		b.cnf.AddToFridgeFromBuyList = true
		b.cnf.UserInsert = true
		msg.Text = fmt.Sprintf("%s\n%s", b.cnf.AddingFridgeListBuy, b.que.GetBuyList(mes.From.UserName, *b.cnf))

	case b.cnf.BuyOrNew.Keyboard[0][1].Text:
		b.cnf.AddToFridge = true
		b.cnf.UserInsert = true
		msg.Text = b.cnf.AddingFridgeListNew

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
		msg.Text = b.que.AddToBuyList(mes.From.UserName, str[0], str[1], str[2], *b.cnf)
		msg.ReplyMarkup = b.cnf.Cancel
		b.cnf.AddToBuyList = false
		b.cnf.UserInsert = false

	case b.cnf.AddToFridge:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = "Неверный ввод"
			break
		}
		msg.Text = b.que.AddProductToFridge(str[0], mes.From.UserName, str[1], *b.cnf)
		msg.ReplyMarkup = b.cnf.BuyOrNew
		b.cnf.AddToFridge = false
		b.cnf.UserInsert = false

	case b.cnf.AddToFridgeFromBuyList:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = "Неверный ввод"
			break
		}
		msg.Text = b.que.AddProductToFridgeFromBuyList(str[0], mes.From.UserName, str[1], *b.cnf)
		msg.ReplyMarkup = b.cnf.BuyOrNew
		b.cnf.AddToFridgeFromBuyList = false
		b.cnf.UserInsert = false

	case b.cnf.OpenProduct:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = b.cnf.ErrorInsert
			break
		}
		msg.Text = b.que.OpenProduct(mes.From.UserName, str[0], str[1], *b.cnf)
		b.cnf.OpenProduct = false
		b.cnf.UserInsert = false

	case b.cnf.ChangeStatus:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = b.cnf.ErrorInsert
			break
		}
		msg.Text = b.que.ChangeStatus(mes.From.UserName, str[0], str[1], *b.cnf)
		b.cnf.ChangeStatus = false
		b.cnf.UserInsert = false

	case b.cnf.GetStatistic:
		str := strings.Split(mes.Text, " ")
		if len(str) != 2 {
			msg.Text = b.cnf.ErrorInsert
			break
		}
		msg.Text = b.que.GetStats(mes.From.UserName, str[0], str[1], *b.cnf)
		b.cnf.GetStatistic = false
		b.cnf.UserInsert = false

	}
	_, err := b.bot.Send(msg)
	return err
}
