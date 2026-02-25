package controller

import (
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	g    *gin.Engine
	repo app.UserRepository
}

func SetupRoutes(repo app.UserRepository) *Controller {
	controller := &Controller{
		g:    gin.Default(),
		repo: repo,
	}

	controller.g.GET("/health", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "working",
		})
	})

	controller.g.POST("/registration", func(c *gin.Context) {
		var req models.Registration_Request
		c.BindJSON(&req)
		val, err := app.Registration_User(req, controller.repo)
		if err != nil {
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

	})

	controller.g.POST("/login", func(c *gin.Context) {
		var req models.Login_Request
		c.BindJSON(&req)
		val, err := app.Login_User(req, controller.repo)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "problem with server",
			})
			return
		}
		if val {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "login succesfull",
			})
			return
		} else {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"message": "user not found",
			})
			return
		}
	})

	return controller
}

func (controller *Controller) Listen(addr string) {
	controller.g.Run(addr)
}
