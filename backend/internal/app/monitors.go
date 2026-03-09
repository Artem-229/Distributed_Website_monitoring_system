package app

import (
	"Distributed_Website_monitoring_system/internal/models"

	"github.com/google/uuid"
)

type MonitorRepository interface {
	GetMonitors(id uuid.UUID) ([]models.Monitor, error)
	AddMonitor(monitor models.Monitor) (bool, error)
	DeleteMonitor(id uuid.UUID) (bool, error)
	GetMonitor(id uuid.UUID) (models.Monitor, error)
}

func GetMonitors(id uuid.UUID, repo MonitorRepository) ([]models.Monitor, error) {
	get, err := repo.GetMonitors(id)
	if err != nil {
		return nil, err
	}

	return get, nil

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

func DeleteMonitor(id uuid.UUID, repo MonitorRepository) (bool, error) {
	ok, err := repo.DeleteMonitor(id)
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
