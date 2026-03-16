package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type TelegramRepo struct {
	DB *sql.DB
}

func (t *TelegramRepo) AddAlert(alert models.Alert) error {
	query := `INSERT INTO alerts (id, monitor_id, url, response_time, alert_type, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := t.DB.Exec(query,
		alert.Id,
		alert.Monitor_Id,
		alert.Url,
		alert.Response_time,
		alert.Alert_type,
		alert.Created_at,
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *TelegramRepo) GetAlertsByID(id uuid.UUID) ([]models.Alert, error) {
	var alerts []models.Alert

	query := `SELECT id, monitor_id, url, response_time, alert_type, created_at FROM alerts WHERE monitor_id = $1`

	rows, err := t.DB.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var alert models.Alert
		rows.Scan(&alert.Id, &alert.Monitor_Id, &alert.Url, &alert.Response_time, &alert.Alert_type, &alert.Created_at)
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
