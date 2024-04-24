package utils

func CalculateIncomeTaxDetailed(income float64) (float64, error) {
	tax := calculateTotalTax(income)
	return tax, nil
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
