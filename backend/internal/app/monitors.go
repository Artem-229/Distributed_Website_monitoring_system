package app

import (
	"Distributed_Website_monitoring_system/internal/models"

	"github.com/google/uuid"
)

type MonitorRepository interface {
	Monitors(id uuid.UUID) ([]models.Monitor, bool)
}

func GetMonitors(user models.User, repo MonitorRepository) ([]models.Monitor, bool) {
	id := user.ID

	get, check := repo.Monitors(id)
	if check == false {
		return nil, false
	}

	return get, true

}

/* func AddMonitor() {

}

func DeleteMonitor() {

}

func GetSpecificMonitor() {

}
*/
