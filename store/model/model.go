package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type AllowanceGorm struct {
	ID            uint    `gorm:"primaryKey"`
	AllowanceType string  `gorm:"type:varchar(255);not null"`
	Amount        float64 `gorm:"type:decimal(18,2);not null"`
}

func InitializeData(db *gorm.DB) error {
	allowances := []AllowanceGorm{
		{AllowanceType: "Personal", Amount: 60000.0},
		{AllowanceType: "Kreceipt", Amount: 0.0},
	}
	for _, allowance := range allowances {
		if allowance.AllowanceType != "Personal" && allowance.AllowanceType != "Kreceipt" {
			return fmt.Errorf("invalid allowance type: %s", allowance.AllowanceType)
		}
		db.FirstOrCreate(&allowance, Allowance{AllowanceType: allowance.AllowanceType})
	}
	return nil
}
