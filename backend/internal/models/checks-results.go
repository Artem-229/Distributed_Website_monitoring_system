package models

import (
	"time"

	"github.com/google/uuid"
)

type Results struct {
	Id            uuid.UUID
	Monitor_id    uuid.UUID
	Time_Interval int
	Checked_at    time.Time
	Responce_time float64
	Status_ok     bool
}
