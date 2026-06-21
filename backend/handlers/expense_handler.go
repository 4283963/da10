package handlers

import (
	"net/http"
	"strconv"
	"time"
	"transfer-tracker/models"
	"transfer-tracker/services"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler() *ExpenseHandler {
	return &ExpenseHandler{
		service: services.NewExpenseService(),
	}
}

type CreateExpenseRequest struct {
	VehicleID   uint               `json:"vehicle_id" binding:"required"`
	ExpenseType models.ExpenseType `json:"expense_type" binding:"required,oneof=window express other"`
	TypeName    string             `json:"type_name" binding:"required"`
	Amount      float64            `json:"amount" binding:"required,gt=0"`
	Description string             `json:"description"`
	ReceiptNo   string             `json:"receipt_no"`
	PayDate     string             `json:"pay_date"`
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payDate *time.Time
	if req.PayDate != "" {
		if t, err := time.Parse("2006-01-02", req.PayDate); err == nil {
			payDate = &t
		}
	}

	expense, err := h.service.CreateExpense(
		req.VehicleID,
		req.ExpenseType,
		req.TypeName,
		req.Amount,
		req.Description,
		req.ReceiptNo,
		payDate,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": expense})
}

func (h *ExpenseHandler) GetStatistics(c *gin.Context) {
	vehicleIDStr := c.Param("vehicleId")
	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的车辆ID"})
		return
	}

	stats, err := h.service.GetStatistics(uint(vehicleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
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
