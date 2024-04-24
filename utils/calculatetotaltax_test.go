package utils

import (
	"github.com/pphee/assessment-tax/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateTotalTax(t *testing.T) {
	tests := []struct {
		name   string
		income float64
		want   float64
	}{
		{"No Tax", 100000, 0},
		{"Lowest Bracket", 300000, (300000 - 150000) * 0.1},
		{"Middle Bracket", 750000, (500000-150000)*0.1 + (750000-500000)*0.15},
		{"High Bracket", 1500000, (500000-150000)*0.1 + (1000000-500000)*0.15 + (1500000-1000000)*0.2},
		{"Highest Bracket", 3000000, (500000-150000)*0.1 + (1000000-500000)*0.15 + (2000000-1000000)*0.2 + (3000000-2000000)*0.35},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTotalTax(tt.income)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateTaxBrackets(t *testing.T) {
	tests := []struct {
		name   string
		income float64
		want   []model.TaxBracket
	}{
		{"No Tax", 100000, []model.TaxBracket{{"0-150,000", 0}, {"150,001-500,000", 0}, {"500,001-1,000,000", 0}, {"1,000,001-2,000,000", 0}, {"2,000,001 ขึ้นไป", 0}}},
		{"Highest Bracket", 3000000, []model.TaxBracket{{"0-150,000", 0}, {"150,001-500,000", (500000 - 150000) * 0.1}, {"500,001-1,000,000", (1000000 - 500000) * 0.15}, {"1,000,001-2,000,000", (2000000 - 1000000) * 0.2}, {"2,000,001 ขึ้นไป", (3000000 - 2000000) * 0.35}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTaxBrackets(tt.income)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateIncomeTaxDetailed(t *testing.T) {
	income := 3000000.0
	expectedTax := (500000-150000)*0.1 + (1000000-500000)*0.15 + (2000000-1000000)*0.2 + (3000000-2000000)*0.35
	expectedBrackets := []model.TaxBracket{
		{"0-150,000", 0},
		{"150,001-500,000", (500000 - 150000) * 0.1},
		{"500,001-1,000,000", (1000000 - 500000) * 0.15},
		{"1,000,001-2,000,000", (2000000 - 1000000) * 0.2},
		{"2,000,001 ขึ้นไป", (3000000 - 2000000) * 0.35},
	}

	tax, taxBrackets := CalculateIncomeTaxDetailed(income)
	assert.Equal(t, expectedTax, tax)
	assert.Equal(t, expectedBrackets, taxBrackets)
}
