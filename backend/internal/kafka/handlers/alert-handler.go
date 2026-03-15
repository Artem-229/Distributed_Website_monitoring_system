package handlers

import (
	"Distributed_Website_monitoring_system/internal/models"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type AlertHandler struct{}

func NewAlertHandler() *AlertHandler {
	return &AlertHandler{}
}

func (A *AlertHandler) HandleMessage(msg []byte, offset kafka.Offset) error {
	var checkEvent models.CheckEvent
	err := json.Unmarshal(msg, &checkEvent)
	if err != nil {
		return err
	}

	if checkEvent.Status_ok == false {
		logrus.Error("ALERT: site is unavaliable")
	}

	if checkEvent.ResponseTime > 1000 {
		logrus.Error("ALERT: the connection is unstable")
	}

	return nil
}
