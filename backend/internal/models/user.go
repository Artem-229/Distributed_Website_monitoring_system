package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Username      string
	Login         string
	Password_Hash string
	Created_at    time.Time
	Telegram_id   int64
}

type RegistrationRequest struct {
	Username    string
	Login       string
	Password    string
	Telegram_id int64
}

type LoginRequest struct {
	Login    string
	Password string
}
