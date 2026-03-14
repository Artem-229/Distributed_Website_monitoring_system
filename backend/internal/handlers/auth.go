package handlers

import (
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repo   app.UserRepository
	secret string
}

func NewAuthHandler(repo app.UserRepository, secret string) *AuthHandler {
	return &AuthHandler{
		repo:   repo,
		secret: secret,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	val, token, err := app.LoginUser(req, h.repo, h.secret)
	if err != nil {
		if token == "" {
			c.IndentedJSON(http.StatusConflict, gin.H{
				"message": "problems with token generating",
			})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "problem with server",
		})
		return
	}
	if val {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "login succesfull",
			"token":   token,
		})
		return
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"message": "user not found",
			"token":   token,
		})
		return
	}
}

func (h *AuthHandler) Registration(c *gin.Context) {
	var req models.RegistrationRequest
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	val, err := app.RegistrationUser(req, h.repo)
	if err != nil {
		if err.Error() == "user already exists" {
			c.IndentedJSON(http.StatusConflict, gin.H{
				"message": "user already exists",
			})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "server problems",
		})
		return
	}
	if val {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "registration completed succesfuly"})
		return
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "problems with registration, try again",
		})
		return
	}
}
