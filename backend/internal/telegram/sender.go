package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(bot *tgbotapi.BotAPI, chat *int64, message string) {
	msg := tgbotapi.NewMessage(*chat, message)
	bot.Send(msg)
}
