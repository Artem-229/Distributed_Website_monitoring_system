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
}

type RegistrationRequest struct {
	Username string
	Login    string
	Password string
}

type LoginRequest struct {
	Login    string
	Password string
}
