package telegram

import (
	"Distributed_Website_monitoring_system/internal/app"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	callbackGetMonitors   = "get_monitors"
	callbackMonitorPrefix = "monitor_"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	repo        app.UserRepository
	sendrepo    app.SendRepository
	monitorrepo app.MonitorRepository
}

func NewBot(token string, repo app.UserRepository, sendrepo app.SendRepository, monitorrepo app.MonitorRepository) *Bot {
	rowbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logrus.Error("Problems while connecting to the bot")
	}
	bot := Bot{
		bot:         rowbot,
		repo:        repo,
		sendrepo:    sendrepo,
		monitorrepo: monitorrepo,
	}

	return &bot
}

func BotHandle(b *Bot) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			chatID := update.CallbackQuery.From.ID
			if data == callbackGetMonitors {
				user, err := b.repo.GetByTelegramID(chatID)
				if err != nil {
					logrus.Error(err)
				}
				monitors, err := b.monitorrepo.GetMonitors(user.ID)
				if err != nil {
					logrus.Error(err)
				}
				var rows [][]tgbotapi.InlineKeyboardButton
				for _, m := range monitors {
					btn := tgbotapi.NewInlineKeyboardButtonData(m.Url, "monitor_"+m.Id.String())
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
				}
				keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
				msg := tgbotapi.NewMessage(chatID, "Choose the monitor: ")
				msg.ReplyMarkup = keyboard
				b.bot.Send(msg)
			} else if strings.HasPrefix(data, callbackMonitorPrefix) {
				monitorID := strings.TrimPrefix(data, callbackMonitorPrefix)
				id, _ := uuid.Parse(monitorID)
				alerts, err := b.sendrepo.GetAlertsByID(id)
				if err != nil {
					logrus.Error(err)
				}
				text := fmt.Sprintf("Alerts for monitor:\n")
				for _, a := range alerts {
					text += fmt.Sprintf("- %s | %s | %.0fms\n", a.Alert_type, a.Created_at.Format("02.01 15:04"), a.Response_time)
				}
				if len(alerts) == 0 {
					text = "No alerts for this monitor"
				}
				msg := tgbotapi.NewMessage(chatID, text)
				b.bot.Send(msg)
			}
		}
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/start" {
			btn := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Мои мониторы", callbackGetMonitors),
				),
			)
			msg := tgbotapi.NewMessage(update.Message.From.ID, fmt.Sprintf("Hello! Here is your id: %d", update.Message.From.ID))
			msg.ReplyMarkup = btn
			b.bot.Send(msg)
		}
	}
}
