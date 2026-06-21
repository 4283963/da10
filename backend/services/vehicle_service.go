package services

import (
	"errors"
	"transfer-tracker/database"
	"transfer-tracker/models"
)

type VehicleService struct{}

func NewVehicleService() *VehicleService {
	return &VehicleService{}
}

func (s *VehicleService) CreateVehicle(vin, plateNumber, ownerName, targetCity, remark string) (*models.Vehicle, error) {
	if vin == "" {
		return nil, errors.New("车架号不能为空")
	}

	var existing models.Vehicle
	if err := database.DB.Where("vin = ?", vin).First(&existing).Error; err == nil {
		return nil, errors.New("该车架号已存在")
	}

	vehicle := models.Vehicle{
		VIN:         vin,
		PlateNumber: plateNumber,
		OwnerName:   ownerName,
		TargetCity:  targetCity,
		Remark:      remark,
	}

	if err := database.DB.Create(&vehicle).Error; err != nil {
		return nil, err
	}

	nodes := models.GetDefaultNodes(vehicle.ID)
	if err := database.DB.Create(&nodes).Error; err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (s *VehicleService) GetByVIN(vin string) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := database.DB.Where("vin = ?", vin).First(&vehicle).Error; err != nil {
		return nil, errors.New("未找到该车架号信息")
	}

	var nodes []models.TransferNode
	database.DB.Where("vehicle_id = ?", vehicle.ID).
		Order("id ASC").
		Find(&nodes)
	vehicle.Nodes = nodes

	var expenses []models.Expense
	database.DB.Where("vehicle_id = ?", vehicle.ID).
		Order("created_at DESC").
		Find(&expenses)
	vehicle.Expenses = expenses

	return &vehicle, nil
}

func (s *VehicleService) List() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := database.DB.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (s *VehicleService) Delete(id uint) error {
	if err := database.DB.Where("vehicle_id = ?", id).Delete(&models.TransferNode{}).Error; err != nil {
		return err
	}
	if err := database.DB.Where("vehicle_id = ?", id).Delete(&models.Expense{}).Error; err != nil {
		return err
	}
	if err := database.DB.Delete(&models.Vehicle{}, id).Error; err != nil {
		return err
	}
	return nil
}
