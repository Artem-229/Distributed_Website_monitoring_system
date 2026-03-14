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
	rowid := c.Param("monitor_id")
	id, _ := uuid.Parse(rowid)
	req, err := h.checkrepo.GetChecks(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "problem with data loading",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message":  "accepted",
		"monitors": req,
	})
}
