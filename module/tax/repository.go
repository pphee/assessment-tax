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
