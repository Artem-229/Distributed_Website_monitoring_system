package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type ChecksRepo struct {
	DB *sql.DB
}

func (r *ChecksRepo) AddResultsMonitors(result models.Results) error {
	query := `INSERT INTO monitor_checks (id, monitor_id, response_time, checked_at, status_ok) VALUES ($1, $2, $3, NOW(), $4)`
	_, err := r.DB.Exec(query, result.Id, result.Monitor_id, result.Responce_time, result.Status_ok)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChecksRepo) GetChecks(id uuid.UUID) ([]models.Results, error) {
	query := `SELECT * FROM monitor_checks WHERE id = $1 ORDER BY checked_at DESC LIMIT = 30`
	var res []models.Results
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return res, err
	}
	var ans models.Results

	for rows.Next() {
		rows.Scan(
			&ans.Id,
			&ans.Monitor_id,
			&ans.Checked_at,
			&ans.Responce_time,
			&ans.Status_ok,
		)
		res = append(res, ans)
	}

	return res, nil
}
