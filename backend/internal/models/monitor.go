package models

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	Id            uuid.UUID `json:"Id"`
	Users_id      uuid.UUID `json:"Users_id"`
	Url           string    `json:"Url"`
	Time_interval int       `json:"Time_interval"`
	Is_active     bool      `json:"Is_active"`
	Created_at    time.Time `json:"Created_at"`
}
