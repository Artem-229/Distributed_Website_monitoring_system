package app

import (
	"Distributed_Website_monitoring_system/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByLogin(login string) (models.User, error)
	Create(user models.User) error
	GetByTelegramID(id int64) (models.User, error)
}

func LoginUser(req models.LoginRequest, repo UserRepository, secret string) (bool, string, error) {
	val, err := repo.GetByLogin(req.Login)
	if err != nil {
		return false, "", nil
	}
	if CheckPasswordHash(val.Password_Hash, req.Password) {
		token, err := GenerateJWTToken(val.ID, secret)
		if err != nil {
			return true, "", nil
		}
		return true, token, nil
	}

	return false, "", nil
}

func RegistrationUser(req models.RegistrationRequest, repo UserRepository) (bool, error) {
	hash, err := HashPassword(req.Password)
	if err != nil {
		return false, err
	}
	newuser := models.User{
		ID:            uuid.New(),
		Username:      req.Username,
		Login:         req.Login,
		Password_Hash: hash,
	}
	err = repo.Create(newuser)
	if err != nil {
		return false, err
	}
	return true, nil
}
