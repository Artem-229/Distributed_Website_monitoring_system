package app

import (
	"Distributed_Website_monitoring_system/internal/models"

	"github.com/google/uuid"
)

type SendRepository interface {
	AddAlert(models.Alert) error
	GetAlertsByID(uuid.UUID) ([]models.Alert, error)
}
