package utils

import "github.com/pphee/assessment-tax/internal/model"

func CalculateIncomeTaxDetailed(income float64) (float64, []model.TaxBracket) {
	tax := calculateTotalTax(income)
	taxBrackets := calculateTaxBrackets(income)
	return tax, taxBrackets
}

func calculateTotalTax(income float64) float64 {
	var tax float64

	if income > 2000000 {
		tax += (income - 2000000) * 0.35
		income = 2000000
	}
	if income > 1000000 {
		tax += (income - 1000000) * 0.2
		income = 1000000
	}
	if income > 500000 {
		tax += (income - 500000) * 0.15
		income = 500000
	}
	if income > 150000 {
		tax += (income - 150000) * 0.1
	}

	return tax
}

func calculateTaxBrackets(income float64) []model.TaxBracket {
	taxBrackets := []model.TaxBracket{
		{Level: "0-150,000", Tax: 0},
		{Level: "150,001-500,000", Tax: 0},
		{Level: "500,001-1,000,000", Tax: 0},
		{Level: "1,000,001-2,000,000", Tax: 0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0},
	}

	if income > 2000000 {
		taxBrackets[4].Tax = (income - 2000000) * 0.35
		income = 2000000
	}
	if income > 1000000 {
		taxBrackets[3].Tax = (income - 1000000) * 0.2
		income = 1000000
	}
	if income > 500000 {
		taxBrackets[2].Tax = (income - 500000) * 0.15
		income = 500000
	}
	if income > 150000 {
		taxBrackets[1].Tax = (income - 150000) * 0.1
	}

	return taxBrackets
}
