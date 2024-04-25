package modelgorm

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
	configs := []struct {
		AllowanceType string
		Amount        float64
	}{
		{"PersonalDefault", 60000},
		{"PersonalMax", 100000},
		{"DonationMax", 100000},
		{"KReceiptDefault", 50000},
		{"KReceiptMax", 100000},
	}

	for _, cfg := range configs {
		allowance := AllowanceGorm{
			AllowanceType: cfg.AllowanceType,
			Amount:        cfg.Amount,
		}
		if err := db.FirstOrCreate(&allowance, AllowanceGorm{AllowanceType: cfg.AllowanceType}).Error; err != nil {
			return fmt.Errorf("failed to initialize data for %s: %v", cfg.AllowanceType, err)
		}
	}

	return nil
}
