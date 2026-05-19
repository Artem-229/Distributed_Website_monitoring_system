package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type ChecksRepo struct {
	DB *sql.DB
}

func (r *ChecksRepo) AddResult(result models.Results) error {
	query := `INSERT INTO monitor_checks (id, monitor_id, time_interval, responce_time, checked_at, status_ok, region)
              VALUES ($1, $2, $3, $4, NOW(), $5, $6)`
	_, err := r.DB.Exec(query, result.Id, result.Monitor_id, result.Time_Interval, result.Responce_time, result.Status_ok, result.Region)
	return err
}

func (r *ChecksRepo) GetChecks(id uuid.UUID) ([]models.Results, error) {
	query := `SELECT id, monitor_id, time_interval, responce_time, checked_at, status_ok, region
              FROM monitor_checks
              WHERE monitor_id = $1
              ORDER BY checked_at DESC
              LIMIT 150`

	var res []models.Results
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var ans models.Results
		rows.Scan(
			&ans.Id,
			&ans.Monitor_id,
			&ans.Time_Interval,
			&ans.Responce_time,
			&ans.Checked_at,
			&ans.Status_ok,
			&ans.Region,
		)
		res = append(res, ans)
	}

	return res, nil
}

func (r *ChecksRepo) GetChecksByRegion(id uuid.UUID) (map[string][]models.Results, error) {
	query := `SELECT id, monitor_id, time_interval, responce_time, checked_at, status_ok, region
              FROM monitor_checks
              WHERE monitor_id = $1
              ORDER BY checked_at DESC
              LIMIT 150`

	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]models.Results)
	for rows.Next() {
		var ans models.Results
		rows.Scan(
			&ans.Id,
			&ans.Monitor_id,
			&ans.Time_Interval,
			&ans.Responce_time,
			&ans.Checked_at,
			&ans.Status_ok,
			&ans.Region,
		)
		result[ans.Region] = append(result[ans.Region], ans)
	}

	return result, nil
}
