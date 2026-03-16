package kafkahandler

import (
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/models"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type AlertHandler struct {
	repo   app.SendRepository
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewAlertHandler(repo app.SendRepository) *AlertHandler {
	return &AlertHandler{repo: repo}
}

func (a *AlertHandler) HandleMessage(msg []byte, offset kafka.Offset) error {
	var checkEvent models.CheckEvent
	err := json.Unmarshal(msg, &checkEvent)
	if err != nil {
		return err
	}

	if checkEvent.Status_ok == false {
		alert := models.Alert{
			Id:            uuid.New(),
			Monitor_Id:    uuid.MustParse(checkEvent.MonitorID),
			Url:           checkEvent.Url,
			Response_time: checkEvent.ResponseTime,
			Alert_type:    "Unavaliable",
			Created_at:    time.Now(),
		}
		a.repo.AddAlert(alert)
	}

	if checkEvent.ResponseTime > 700 {
		alert := models.Alert{
			Id:            uuid.New(),
			Monitor_Id:    uuid.MustParse(checkEvent.MonitorID),
			Url:           checkEvent.Url,
			Response_time: checkEvent.ResponseTime,
			Alert_type:    "Slow connection",
			Created_at:    time.Now(),
		}
		a.repo.AddAlert(alert)
	}

	return nil
}
