package app

import (
	"Distributed_Website_monitoring_system/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByLogin(login string) (models.User, error)
	Create(user models.User) error
}

func Login_User(req models.Login_Request, repo UserRepository, secret string) (bool, string, error) {
	val, err := repo.GetByLogin(req.Login)
	if err != nil {
		return false, "", nil
	}
	if Check_Password_Hash(val.Password_Hash, req.Password) {
		token, err := GenerateJWTToken(val.ID, secret)
		if err != nil {
			return true, "", nil
		}
		return true, token, nil
	}

	return false, "", nil
}

func Registration_User(req models.Registration_Request, repo UserRepository) (bool, error) {
	hash, err := Hash_password(req.Password)
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
