package tax

import (
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/pphee/assessment-tax/internal/model"
	"github.com/pphee/assessment-tax/utils"
	"io"
	"log"
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

	personalDefault := allowances[0].Amount
	personalMax := allowances[1].Amount
	donationMax := allowances[2].Amount
	kReceiptDefault := allowances[3].Amount
	kReceiptMax := allowances[4].Amount

	var totalDeductions float64
	for _, allowance := range req.Allowances {
		if allowance.Amount < 0 {
			return 0, nil, errors.New("allowance amount cannot be negative")
		}

		switch allowance.AllowanceType {
		case "personal":
			if allowance.Amount > personalMax || allowance.Amount < 10000 {
				return 0, nil, fmt.Errorf("personal allowance amount must be between 10000 and %f", personalMax)
			}
		case "donation":
			if allowance.Amount > donationMax {
				allowance.Amount = donationMax
			}
		case "k-receipt":
			if allowance.Amount > 0 {
				allowance.Amount = kReceiptDefault
			}
		case "k-receipt-admin":
			if allowance.Amount > kReceiptMax {
				allowance.Amount = kReceiptMax
			}
		}
		totalDeductions += allowance.Amount
	}

	if req.WHT < 0 || req.WHT > req.TotalIncome {
		return 0, nil, errors.New("invalid WHT value")
	}

	taxableIncome := req.TotalIncome - totalDeductions - personalDefault
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

func (service *TaxService) SetPersonalDeduction(amount float64) error {
	if amount < 10000 || amount > 100000 {
		return errors.New("amount must be between 10,000 and 100,000")
	}

	err := service.Repo.SetPersonalDeduction(amount)
	if err != nil {
		return fmt.Errorf("failed to set personal deduction: %w", err)
	}

	return nil
}

func (service *TaxService) TaxFromFile(file io.Reader) ([]model.TotalIncomeCsv, error) {
	var totalIncomeCsv []model.TotalIncomeCsv
	if err := gocsv.Unmarshal(file, &totalIncomeCsv); err != nil {
		log.Fatal("Failed to unmarshal CSV: ", err)
		return nil, err
	}
	return totalIncomeCsv, nil
}

func (service *TaxService) SetKReceiptDeduction(amount float64) error {
	if amount < 0 || amount > 100000 {
		return errors.New("amount must be between 0 and 100,000")
	}

	err := service.Repo.SetKreceiptDeduction(amount)
	if err != nil {
		return fmt.Errorf("failed to set K-receipt deduction: %w", err)
	}

	return nil
}
