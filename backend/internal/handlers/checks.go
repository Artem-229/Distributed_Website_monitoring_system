package handlers

import (
	"Distributed_Website_monitoring_system/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CheckHandler struct {
	checkrepo app.ChecksRepository
}

func NewCheckHandler(checkrepo app.ChecksRepository) *CheckHandler {
	return &CheckHandler{
		checkrepo: checkrepo,
	}
}

func (h *CheckHandler) Check(c *gin.Context) {
	id, err := uuid.Parse(c.Param("monitor_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid monitor_id"})
		return
	}

	req, err := h.checkrepo.GetChecks(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "problem with data loading"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "accepted",
		"monitors": req,
	})
}

func (h *CheckHandler) ChecksByRegion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("monitor_id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid monitor_id"})
		return
	}

	req, err := h.checkrepo.GetChecksByRegion(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "problem with data loading"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "accepted",
		"regions": req,
	})
}
