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

func (service *TaxService) CalculateTax(req model.TaxRequest) (float64, []model.TaxBracket, error) {
	allowances, err := service.Repo.GetAllowanceConfig()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to retrieve allowance configuration: %w", err)
	}

	var totalDeductions float64
	for _, allowance := range req.Allowances {
		if allowance.Amount < 0 {
			return 0, nil, errors.New("allowance amount cannot be negative")
		}
		PersonalMax := allowances[1].Amount
		DonationMax := allowances[2].Amount
		switch allowance.AllowanceType {
		case "personal":
			if allowance.Amount > PersonalMax || allowance.Amount < 10000 {
				return 0, nil, fmt.Errorf("personal allowance amount must be between 10000 and %f", PersonalMax)
			}
		case "donation":
			if allowance.Amount > DonationMax {
				allowance.Amount = DonationMax
			}
		}
		totalDeductions += allowance.Amount
	}

	if req.WHT < 0 || req.WHT > req.TotalIncome {
		return 0, nil, errors.New("invalid WHT value")
	}
	PersonalDefault := allowances[0].Amount

	taxableIncome := req.TotalIncome - totalDeductions - PersonalDefault
	tax, taxBrackets := utils.CalculateIncomeTaxDetailed(taxableIncome)
	tax -= req.WHT
	for i := range taxBrackets {
		taxBrackets[i].Tax -= req.WHT
		if taxBrackets[i].Tax < 0 {
			taxBrackets[i].Tax = 0
		}
	}

	return tax, taxBrackets, nil
}
