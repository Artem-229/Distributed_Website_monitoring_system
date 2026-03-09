package app

import (
	"Distributed_Website_monitoring_system/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ChecksRepository interface {
	AddResultsMonitors(models.Results) error
	GetChecks(id uuid.UUID) ([]models.Results, error)
}

func CheckPing(monitor models.Monitor, repo ChecksRepository) (float64, bool, error) {

	start := time.Now()

	req, err := http.Get(monitor.Url)
	if err != nil {
		return 0, false, err
	}
	if req.StatusCode >= 500 {
		return 0, true, err
	}

	end := time.Since(start)

	res := models.Results{
		Id:            uuid.New(),
		Monitor_id:    monitor.Id,
		Time_Interval: monitor.Time_interval,
		Checked_at:    start,
		Status_ok:     true,
		Responce_time: float64(end.Milliseconds()),
	}

	err = repo.AddResultsMonitors(res)

	if err != nil {
		return float64(end.Milliseconds()), true, err
	}

	return float64(end.Milliseconds()), true, nil

}

func GetChecks(id uuid.UUID, repo ChecksRepository) ([]models.Results, error) {
	res, err := repo.GetChecks(id)
	if err != nil {
		return res, err
	}

	return res, nil
}
