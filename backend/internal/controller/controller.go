package controller

import (
	"Distributed_Website_monitoring_system/internal/handlers"
	"Distributed_Website_monitoring_system/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	g *gin.Engine
}

func SetupRoutes(auth *handlers.AuthHandler, monitor *handlers.MonitorHandler, check *handlers.CheckHandler, health *handlers.HealthHandler, secret string) *Controller {
	controller := &Controller{
		g: gin.Default(),
	}

	controller.g.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authorized := controller.g.Group("/api", middleware.CheckJWT(secret))
	controller.g.GET("/health", health.HealthCheck)
	controller.g.POST("/registration", auth.Registration)
	controller.g.POST("/login", auth.Login)
	authorized.GET("/monitors", monitor.GetMonitors)
	authorized.POST("/addmonitor", monitor.AddMonitor)
	authorized.POST("/deletemonitor", monitor.DeleteMonitor)
	authorized.POST("/getmonitor", monitor.GetMonitor)
	authorized.POST("/checks/:monitor_id", check.Check)

	return controller
}

func (controller *Controller) Listen(addr string) {
	controller.g.Run(addr)
}
