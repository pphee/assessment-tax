package tax

import (
	"bytes"
	"github.com/pphee/assessment-tax/internal/model"
	modelgorm "github.com/pphee/assessment-tax/store/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetAllowanceConfig() ([]modelgorm.AllowanceGorm, error) {
	args := m.Called()
	return args.Get(0).([]modelgorm.AllowanceGorm), args.Error(1)
}

func (m *MockRepo) SetPersonalDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *MockRepo) SetKreceiptDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func TestCalculateTax(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewTaxService(mockRepo)

	allowances := []modelgorm.AllowanceGorm{
		{Amount: 60000.00},  // personalDefault
		{Amount: 100000.00}, // personalMax
		{Amount: 100000.00}, // donationMax
		{Amount: 50000.00},  // kReceiptDefault
		{Amount: 100000.00}, // kReceiptMax
	}

	mockRepo.On("GetAllowanceConfig").Return(allowances, nil)

	req := model.TaxRequest{
		TotalIncome: 500000.0,
		WHT:         25000.0,
		Allowances: []model.Allowance{
			{AllowanceType: "donation", Amount: 0},
		},
	}

	expectedTax := 4000.0
	expectedBrackets := []model.TaxBracket{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 4000.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	tax, taxBrackets, err := service.CalculateTax(req)

	assert.Nil(t, err)
	assert.Equal(t, expectedTax, tax)
	assert.ElementsMatch(t, expectedBrackets, taxBrackets)
}

func TestSettPersonalDeduction(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewTaxService(mockRepo)
	amount := 50000.0

	mockRepo.On("SetPersonalDeduction", amount).Return(nil)

	err := service.SetPersonalDeduction(amount)
	assert.Nil(t, err)
}

func TestSetKReceiptDeduction(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewTaxService(mockRepo)
	amount := 40000.0

	mockRepo.On("SetKreceiptDeduction", amount).Return(nil)

	err := service.SetKReceiptDeduction(amount)
	assert.Nil(t, err)
}

func TestTaxFromFileSuccess(t *testing.T) {
	// Sample CSV content with headers matching the struct tags in TotalIncomeCsv
	csvContent := `totalIncome,wht,donation
50000,5000,200
60000,6000,300`
	reader := bytes.NewBufferString(csvContent)

	service := TaxService{}

	// Call the method
	result, err := service.TaxFromFile(reader)

	// Define the expected output
	expected := []model.TotalIncomeCsv{
		{TotalIncome: 50000, WHT: 5000, Donation: 200},
		{TotalIncome: 60000, WHT: 6000, Donation: 300},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestTaxFromFileError(t *testing.T) {
	csvContent := `totalIncome,wht,donation
50000,notanumber,200`
	reader := bytes.NewBufferString(csvContent)

	service := TaxService{}

	_, err := service.TaxFromFile(reader)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "strconv.ParseFloat: parsing \"notanumber\"")
}
