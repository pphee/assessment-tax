package tax

import (
	"github.com/pphee/assessment-tax/store/model"
	"gorm.io/gorm"
)

type TaxRepository struct {
	DB *gorm.DB
}

func NewTaxRepository(db *gorm.DB) *TaxRepository {
	return &TaxRepository{DB: db}
}

func (repo *TaxRepository) GetAllowanceConfig() ([]model.AllowanceGorm, error) {
	var allowances []model.AllowanceGorm
	err := repo.DB.Find(&allowances).Error
	if err != nil {
		return nil, err
	}
	return allowances, nil
}

func (repo *TaxRepository) SetPersonalDeduction(amount float64) error {
	allowance := model.AllowanceGorm{
		AllowanceType: "personalDeduction",
		Amount:        amount,
	}
	if err := repo.DB.Save(&allowance).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TaxRepository) SetKreceiptDeduction(amount float64) error {
	allowance := model.AllowanceGorm{
		AllowanceType: "kReceipt",
		Amount:        amount,
	}
	if err := repo.DB.Save(&allowance).Error; err != nil {
		return err
	}
	return nil
}
