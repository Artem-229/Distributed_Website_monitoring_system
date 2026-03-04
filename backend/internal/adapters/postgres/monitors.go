package postgres

import (
	"Distributed_Website_monitoring_system/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

type MonitorRepo struct {
	DB *sql.DB
}

func (r *MonitorRepo) Monitors(id uuid.UUID) ([]models.Monitor, bool) {
	query := `
		SELECT *
		FROM monitors
		WHERE userid = $1
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
