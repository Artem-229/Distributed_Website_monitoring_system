package app

import (
	"Distributed_Website_monitoring_system/internal/kafka/producer"
	"Distributed_Website_monitoring_system/internal/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ChecksRepository interface {
	AddResult(models.Results) error
	GetChecks(id uuid.UUID) ([]models.Results, error)
	GetChecksByRegion(id uuid.UUID) (map[string][]models.Results, error)
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func Ping(url string) (float64, bool, error) {
	start := time.Now()
	resp, err := httpClient.Get(url)
	elapsed := time.Since(start)

	statusOk := err == nil && resp != nil && resp.StatusCode < 500
	if resp != nil {
		resp.Body.Close()
	}
	return float64(elapsed.Milliseconds()), statusOk, err
}

func SaveCheck(monitor models.Monitor, repo ChecksRepository, region string, responseTime float64, statusOk bool) error {
	res := models.Results{
		Id:            uuid.New(),
		Monitor_id:    monitor.Id,
		Time_Interval: monitor.Time_interval,
		Checked_at:    time.Now(),
		Status_ok:     statusOk,
		Responce_time: responseTime,
		Region:        region,
	}
	return repo.AddResult(res)
}

func SendKafkaEvent(prod *producer.Producer, monitor models.Monitor, responseTime float64, statusOk bool) error {
	if prod == nil {
		return nil
	}
	checkEvent := models.CheckEvent{
		MonitorID:    monitor.Id.String(),
		Url:          monitor.Url,
		Status_ok:    statusOk,
		ResponseTime: responseTime,
	}
	data, _ := json.Marshal(checkEvent)
	return prod.Produce(string(data), "monitor.results")
}

func GetChecks(id uuid.UUID, repo ChecksRepository) ([]models.Results, error) {
	return repo.GetChecks(id)
}
