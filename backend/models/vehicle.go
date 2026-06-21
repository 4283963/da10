package models

import "time"

type Vehicle struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	VIN         string         `gorm:"uniqueIndex;size:50;not null" json:"vin"`
	PlateNumber string         `gorm:"size:20" json:"plate_number"`
	OwnerName   string         `gorm:"size:50" json:"owner_name"`
	TargetCity  string         `gorm:"size:50" json:"target_city"`
	Remark      string         `gorm:"size:255" json:"remark"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Nodes       []TransferNode `gorm:"foreignKey:VehicleID" json:"nodes,omitempty"`
	Expenses    []Expense      `gorm:"foreignKey:VehicleID" json:"expenses,omitempty"`
}
