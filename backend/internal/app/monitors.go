package app

import (
	"Distributed_Website_monitoring_system/internal/models"

	"github.com/google/uuid"
)

type MonitorRepository interface {
	GetMonitors(id uuid.UUID) ([]models.Monitor, bool)
	AddMonitor(monitor models.Monitor) (bool, error)
	DeleteMonitor(monitor models.Monitor) (bool, error)
	GetMonitor(id uuid.UUID) (models.Monitor, error)
}

func GetMonitors(id uuid.UUID, repo MonitorRepository) ([]models.Monitor, bool) {
	get, check := repo.GetMonitors(id)
	if check == false {
		return nil, false
	}

	return get, true

}

func AddMonitor(monitor models.Monitor, repo MonitorRepository) (bool, error) {
	ok, err := repo.AddMonitor(monitor)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	return false, nil
}

func DeleteMonitor(monitor models.Monitor, repo MonitorRepository) (bool, error) {
	ok, err := repo.DeleteMonitor(monitor)
	if err != nil {
		return false, err
	}

	if ok {
		return true, err
	}

	return false, nil
}

func GetMonitor(id uuid.UUID, repo MonitorRepository) (models.Monitor, error) {
	mon, err := repo.GetMonitor(id)
	if err != nil {
		return mon, err
	}

	return mon, nil
}
