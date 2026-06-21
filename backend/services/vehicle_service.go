package services

import (
	"errors"
	"transfer-tracker/database"
	"transfer-tracker/models"

	"gorm.io/gorm"
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
	if err := database.DB.Where("vin = ?", vin).
		Preload("Nodes", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Expenses", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&vehicle).Error; err != nil {
		return nil, errors.New("未找到该车架号信息")
	}
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
