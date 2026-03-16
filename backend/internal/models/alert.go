package models

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	Id            uuid.UUID
	Monitor_Id    uuid.UUID
	Url           string
	Response_time float64
	Alert_type    string
	Created_at    time.Time
}
