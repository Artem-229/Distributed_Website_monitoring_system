package models

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	Id            uuid.UUID
	Users_id      uuid.UUID
	Url           string
	Time_interval int
	Is_active     bool
	Created_at    time.Time
}
