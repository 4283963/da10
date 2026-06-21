package services

import (
	"errors"
	"time"
	"transfer-tracker/database"
	"transfer-tracker/models"
)

type ExpenseService struct{}

func NewExpenseService() *ExpenseService {
	return &ExpenseService{}
}

func (s *ExpenseService) CreateExpense(vehicleID uint, expenseType models.ExpenseType, typeName string, amount float64, description, receiptNo string, payDate *time.Time) (*models.Expense, error) {
	if vehicleID == 0 {
		return nil, errors.New("车辆ID不能为空")
	}
	if typeName == "" {
		return nil, errors.New("费用名称不能为空")
	}
	if amount <= 0 {
		return nil, errors.New("费用金额必须大于0")
	}

	expense := models.Expense{
		VehicleID:   vehicleID,
		ExpenseType: expenseType,
		TypeName:    typeName,
		Amount:      amount,
		Description: description,
		ReceiptNo:   receiptNo,
		PayDate:     payDate,
	}

	if err := database.DB.Create(&expense).Error; err != nil {
		return nil, err
	}

	return &expense, nil
}

func (s *ExpenseService) GetByVehicleID(vehicleID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	if err := database.DB.Where("vehicle_id = ?", vehicleID).Order("created_at DESC").Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (s *ExpenseService) GetStatistics(vehicleID uint) (map[string]interface{}, error) {
	type StatResult struct {
		ExpenseType string  `gorm:"column:expense_type"`
		TotalAmount float64 `gorm:"column:total_amount"`
		Count       int     `gorm:"column:count"`
	}

	var stats []StatResult
	database.DB.Model(&models.Expense{}).
		Select("expense_type, SUM(amount) as total_amount, COUNT(*) as count").
		Where("vehicle_id = ?", vehicleID).
		Group("expense_type").
		Scan(&stats)

	var totalAmount float64
	var windowAmount float64
	var expressAmount float64
	var otherAmount float64
	var totalCount int

	for _, s := range stats {
		totalAmount += s.TotalAmount
		totalCount += s.Count
		switch models.ExpenseType(s.ExpenseType) {
		case models.ExpenseWindow:
			windowAmount = s.TotalAmount
		case models.ExpenseExpress:
			expressAmount = s.TotalAmount
		default:
			otherAmount = s.TotalAmount
		}
	}

	var expenses []models.Expense
	database.DB.Where("vehicle_id = ?", vehicleID).
		Order("created_at DESC").
		Find(&expenses)

	return map[string]interface{}{
		"total_amount":   totalAmount,
		"window_amount":  windowAmount,
		"express_amount": expressAmount,
		"other_amount":   otherAmount,
		"total_count":    totalCount,
		"expenses":       expenses,
	}, nil
}

func (s *ExpenseService) Delete(id uint) error {
	if err := database.DB.Delete(&models.Expense{}, id).Error; err != nil {
		return err
	}
	return nil
}
