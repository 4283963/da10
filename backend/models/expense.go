package models

import "time"

type ExpenseType string

const (
	ExpenseWindow  ExpenseType = "window"
	ExpenseExpress ExpenseType = "express"
	ExpenseOther   ExpenseType = "other"
)

type Expense struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	VehicleID   uint        `gorm:"index;not null" json:"vehicle_id"`
	ExpenseType ExpenseType `gorm:"size:20;not null" json:"expense_type"`
	TypeName    string      `gorm:"size:50;not null" json:"type_name"`
	Amount      float64     `gorm:"type:decimal(10,2);not null" json:"amount"`
	Description string      `gorm:"size:255" json:"description"`
	PayDate     *time.Time  `json:"pay_date"`
	ReceiptNo   string      `gorm:"size:100" json:"receipt_no"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
