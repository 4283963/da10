package handlers

import (
	"net/http"
	"strconv"
	"transfer-tracker/services"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	service *services.VehicleService
}

func NewVehicleHandler() *VehicleHandler {
	return &VehicleHandler{
		service: services.NewVehicleService(),
	}
}

type CreateVehicleRequest struct {
	VIN         string `json:"vin" binding:"required"`
	PlateNumber string `json:"plate_number"`
	OwnerName   string `json:"owner_name"`
	TargetCity  string `json:"target_city"`
	Remark      string `json:"remark"`
}

func (h *VehicleHandler) Create(c *gin.Context) {
	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle, err := h.service.CreateVehicle(req.VIN, req.PlateNumber, req.OwnerName, req.TargetCity, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": vehicle})
}

func (h *VehicleHandler) GetByVIN(c *gin.Context) {
	vin := c.Param("vin")
	if vin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "车架号不能为空"})
		return
	}

	vehicle, err := h.service.GetByVIN(vin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vehicle})
}

func (h *VehicleHandler) List(c *gin.Context) {
	vehicles, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vehicles})
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
