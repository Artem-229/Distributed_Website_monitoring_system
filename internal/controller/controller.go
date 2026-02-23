package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	g *gin.Engine
}

func SetupRoutes() *Controller {
	controller := &Controller{
		g: gin.Default(),
	}

	controller.g.GET("/check", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "working",
		})
	})

	return controller
}

func (controller *Controller) Listen(addr string) {
	controller.g.Run(addr)
}
