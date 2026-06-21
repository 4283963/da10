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
	expenses, err := s.GetByVehicleID(vehicleID)
	if err != nil {
		return nil, err
	}

	var totalAmount float64
	var windowAmount float64
	var expressAmount float64
	var otherAmount float64

	for _, e := range expenses {
		totalAmount += e.Amount
		switch e.ExpenseType {
		case models.ExpenseWindow:
			windowAmount += e.Amount
		case models.ExpenseExpress:
			expressAmount += e.Amount
		default:
			otherAmount += e.Amount
		}
	}

	return map[string]interface{}{
		"total_amount":   totalAmount,
		"window_amount":  windowAmount,
		"express_amount": expressAmount,
		"other_amount":   otherAmount,
		"total_count":    len(expenses),
		"expenses":       expenses,
	}, nil
}

func (s *ExpenseService) Delete(id uint) error {
	if err := database.DB.Delete(&models.Expense{}, id).Error; err != nil {
		return err
	}
	return nil
}
