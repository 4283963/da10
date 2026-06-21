package handlers

import (
	"net/http"
	"strconv"
	"transfer-tracker/models"
	"transfer-tracker/services"

	"github.com/gin-gonic/gin"
)

type TransferNodeHandler struct {
	service *services.TransferNodeService
}

func NewTransferNodeHandler() *TransferNodeHandler {
	return &TransferNodeHandler{
		service: services.NewTransferNodeService(),
	}
}

type UpdateNodeRequest struct {
	Status      models.NodeStatus `json:"status" binding:"required,oneof=pending doing done"`
	Description string            `json:"description"`
	Operator    string            `json:"operator"`
	Location    string            `json:"location"`
}

func (h *TransferNodeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node, err := h.service.UpdateNode(uint(id), req.Status, req.Description, req.Operator, req.Location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": node})
}

func (h *TransferNodeHandler) GetProgress(c *gin.Context) {
	vehicleIDStr := c.Param("vehicleId")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的车辆ID"})
		return
	}

	progress, err := h.service.GetProgress(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}
