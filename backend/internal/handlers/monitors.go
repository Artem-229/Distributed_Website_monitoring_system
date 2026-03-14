package handlers

import (
	"Distributed_Website_monitoring_system/internal/app"
	"Distributed_Website_monitoring_system/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MonitorHandler struct {
	monitorrepo app.MonitorRepository
}

func NewMonitorHandler(monitorrepo app.MonitorRepository) *MonitorHandler {
	return &MonitorHandler{
		monitorrepo: monitorrepo,
	}
}

func (m *MonitorHandler) GetMonitors(c *gin.Context) {
	UserIDraw, _ := c.Get("UserID")
	UserIDstring := UserIDraw.(string)

	UserID, err := uuid.Parse(UserIDstring)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "something wrong with id",
		})
		return
	}

	ans, err := m.monitorrepo.GetMonitors(UserID)
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
}

func (m *MonitorHandler) AddMonitor(c *gin.Context) {
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

	ok, err := app.AddMonitor(req, m.monitorrepo)
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
}

func (m *MonitorHandler) DeleteMonitor(c *gin.Context) {
	var mon models.Monitor
	err := c.BindJSON(&mon)

	ok, err := app.DeleteMonitor(mon.Id, m.monitorrepo)
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
}

func (m *MonitorHandler) GetMonitor(c *gin.Context) {
	UserIDRaw, _ := c.Get("UserID")
	UserIDstring := UserIDRaw.(string)
	UserID, _ := uuid.Parse(UserIDstring)

	mon, err := m.monitorrepo.GetMonitor(UserID)
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
}
