package controller

import (
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/middleware"
	"Distributed_Website_monitoring_system/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	g       *gin.Engine
	repo    app.UserRepository
	secret  string
	monitor app.MonitorRepository
	checks  app.ChecksRepository
}

func SetupRoutes(repo app.UserRepository, secret string, monitor app.MonitorRepository, checks app.ChecksRepository) *Controller {
	controller := &Controller{
		g:       gin.Default(),
		repo:    repo,
		secret:  secret,
		monitor: monitor,
		checks:  checks,
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

	controller.g.GET("/health", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "working",
		})
	})

	controller.g.POST("/registration", func(c *gin.Context) {
		var req models.Registration_Request
		err := c.BindJSON(&req)
		if err != nil {
			return
		}
		val, err := app.Registration_User(req, controller.repo)
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

	})

	controller.g.POST("/login", func(c *gin.Context) {
		var req models.Login_Request
		err := c.BindJSON(&req)
		if err != nil {
			return
		}
		val, token, err := app.Login_User(req, controller.repo, controller.secret)
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
	})

	authorized.GET("/monitors", func(c *gin.Context) {

		UserIDraw, _ := c.Get("UserID")
		UserIDstring := UserIDraw.(string)

		UserID, err := uuid.Parse(UserIDstring)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "something wrong with id",
			})
			return
		}

		ans, err := monitor.GetMonitors(UserID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "problems with monitors parsing",
			})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"message":  "accepted",
			"monitors": ans,
		})
	})

	authorized.POST("/addmonitor", func(c *gin.Context) {
		var req models.Monitor
		err := c.BindJSON(&req)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "problem with data",
			})
			return
		}

		req.Id = uuid.New()
		UserIDRaw, _ := c.Get("UserID")
		UserIDstring := UserIDRaw.(string)
		req.Users_id, _ = uuid.Parse(UserIDstring)

		ok, err := app.AddMonitor(req, monitor)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "problems with adding a monitor",
			})
			return

		}
		if ok {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "monitor added succesfully",
			})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server problems",
			})
		}

	})

	authorized.POST("/deletemonitor", func(c *gin.Context) {
		var mon models.Monitor
		err := c.BindJSON(&mon)

		ok, err := app.DeleteMonitor(mon.Id, monitor)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "error while loading the monitor",
			})
			return
		}
		if ok {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": "accepted",
			})
		}
	})

	authorized.POST("/getmonitor", func(c *gin.Context) {
		UserIDRaw, _ := c.Get("UserID")
		UserIDstring := UserIDRaw.(string)
		UserID, _ := uuid.Parse(UserIDstring)

		mon, err := monitor.GetMonitor(UserID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "couldnt get the monitor",
			})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "accepted",
			"monitor": mon,
		})
	})

	authorized.POST("/checks/:monitor_id", func(c *gin.Context) {
		rowid := c.Param("monitor_id")
		id, _ := uuid.Parse(rowid)
		req, err := checks.GetChecks(id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "problem with data loading",
			})
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"message":  "accepted",
			"monitors": req,
		})

	})

	return controller
}

func (controller *Controller) Listen(addr string) {
	controller.g.Run(addr)
}
