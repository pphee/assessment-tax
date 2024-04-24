package tax

import (
	"errors"
	"fmt"
	"github.com/pphee/assessment-tax/internal/model"
	"github.com/pphee/assessment-tax/utils"
)

type TaxService struct {
	Repo *TaxRepository
}

func NewTaxService(repo *TaxRepository) *TaxService {
	return &TaxService{Repo: repo}
}

func (service *TaxService) CalculateTax(req model.TaxRequest) (float64, error) {
	allowances, err := service.Repo.GetAllowanceConfig()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve allowance configuration: %w", err)
	}

	var totalDeductions float64
	for _, allowance := range req.Allowances {
		if allowance.Amount < 0 {
			return 0, errors.New("allowance amount cannot be negative")
		}
		switch allowance.AllowanceType {
		case "personal":
			if allowance.Amount > allowances[1].Amount || allowance.Amount < 10000 {
				return 0, fmt.Errorf("personal allowance amount must be between 10000 and %f", allowances[1].Amount)
			}
		case "donation":
			if allowance.Amount > allowances[2].Amount {
				allowance.Amount = allowances[2].Amount
			}
		case "k-receipt":
			if allowance.Amount > allowances[3].Amount {
				allowance.Amount = allowances[3].Amount
			}
		}
		totalDeductions += allowance.Amount
	}

	if req.WHT < 0 || req.WHT > req.TotalIncome {
		return 0, errors.New("invalid WHT value")
	}

	taxableIncome := req.TotalIncome - totalDeductions - allowances[0].Amount
	tax, err := utils.CalculateIncomeTaxDetailed(taxableIncome)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate income tax: %w", err)
	}
	tax -= req.WHT

	return tax, nil
}
