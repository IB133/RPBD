package tgbot

import (
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(mes *tgbotapi.Message) error {
	switch mes.Command() {
	case "start":
		msg := tgbotapi.NewMessage(mes.Chat.ID, b.messages.Start)
		msg.ReplyMarkup = b.keyboards.Main
		_, err := b.bot.Send(msg)
		return err
	default:
		msg := tgbotapi.NewMessage(mes.Chat.ID, "Неверная команда")
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleMessages(mes *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(mes.Chat.ID, mes.Text)
	switch mes.Text {
	case b.keyboards.Main.Keyboard[0][0].Text:
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		handleAddToBuyList(&msg)

	case b.keyboards.Main.Keyboard[0][1].Text:
		msg.ReplyToMessageID = mes.MessageID

	case b.keyboards.Main.Keyboard[1][0].Text:
		msg.ReplyToMessageID = mes.MessageID

	case b.keyboards.Main.Keyboard[1][1].Text:
		msg.ReplyToMessageID = mes.MessageID

	case b.keyboards.Main.Keyboard[1][2].Text:
		msg.ReplyToMessageID = mes.MessageID

	case b.keyboards.Main.Keyboard[2][0].Text:
		msg.ReplyToMessageID = mes.MessageID

	case b.keyboards.Main.Keyboard[2][1].Text:
		msg.ReplyToMessageID = mes.MessageID

	default:
		msg := tgbotapi.NewMessage(mes.Chat.ID, "Выберите пункт из меню")
		msg.ReplyMarkup = b.keyboards.Main
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleIncorrect(mes *tgbotapi.Message) error {
	if reflect.TypeOf(mes.Text).Kind() != reflect.String && mes.Text == "" {
		msg := tgbotapi.NewMessage(mes.Chat.ID, "Хуйню ввел")
		_, err := b.bot.Send(msg)
		return err
	}
	return nil
}

func handleAddToBuyList(msg *tgbotapi.MessageConfig) error {
	msg.Text = "Введите вес в кг"
	msg.Text = "Введите aboba в кг"
	return nil
}
