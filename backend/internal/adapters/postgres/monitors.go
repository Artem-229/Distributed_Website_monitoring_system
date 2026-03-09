package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type MonitorRepo struct {
	DB *sql.DB
}

func (r *MonitorRepo) GetMonitors(id uuid.UUID) ([]models.Monitor, bool) {
	query := `
		SELECT *
		FROM monitors
		WHERE user_id = $1
	`

	res := make([]models.Monitor, 0)

	var ans models.Monitor

	rows, err := r.DB.Query(query, id)
	if err != nil {
		return res, false
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&ans.Id,
			&ans.Url,
			&ans.Time_interval,
			&ans.Is_active,
			&ans.Created_at,
		)
		res = append(res, ans)
	}

	return res, true
}

func (r *MonitorRepo) AddMonitor(monitor models.Monitor) (bool, error) {
	query := `INSERT INTO monitors (id, user_id, url, time_interval, is_active, created_at) VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err := r.DB.Exec(query, monitor.Id, monitor.Users_id, monitor.Url, monitor.Time_interval, monitor.Is_active)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (r *MonitorRepo) DeleteMonitor(monitor models.Monitor) (bool, error) {
	query := `
		DELETE 
		FROM monitors 
		WHERE id = $1
	`

	_, err := r.DB.Exec(query, monitor.Id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *MonitorRepo) GetMonitor(id uuid.UUID) (models.Monitor, error) {
	query := `
		SELECT * 
		FROM monitors 
		WHERE id = $1
	`

	var req models.Monitor
	mon := r.DB.QueryRow(query, id)

	err := mon.Scan(&req.Id, &req.Url, &req.Time_interval, &req.Is_active, &req.Created_at)
	if err != nil {
		return req, err
	}

	return req, nil
}
