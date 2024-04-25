package tax

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/pphee/assessment-tax/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockTaxService struct {
	mock.Mock
}

func (m *MockTaxService) CalculateTax(req model.TaxRequest) (float64, []model.TaxBracket, error) {
	args := m.Called(req)
	return args.Get(0).(float64), args.Get(1).([]model.TaxBracket), args.Error(2)
}

func (m *MockTaxService) SetPersonalDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *MockTaxService) TaxFromFile(file io.Reader) ([]model.TotalIncomeCsv, error) {
	args := m.Called(file)
	return args.Get(0).([]model.TotalIncomeCsv), args.Error(1)
}

func (m *MockTaxService) SetKReceiptDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func TestTaxHandler_PostTaxCalculation(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"totalIncome": 500000, "wht": 25000, "allowances":[{"allowanceType":"k-receipt","amount":50000}]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mocking the TaxService
	mockTaxService := new(MockTaxService)
	mockTaxService.On("CalculateTax", mock.Anything).Return(12345.67, []model.TaxBracket{{Level: "Low", Tax: 1000}}, nil)

	h := &TaxHandler{
		TaxService: mockTaxService, // Inject the mocked service
	}

	// Perform the test
	if assert.NoError(t, h.PostTaxCalculation(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestTaxHandler_SetPersonalDeduction(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 50000}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mocking the TaxService
	mockTaxService := new(MockTaxService)
	mockTaxService.On("SetPersonalDeduction", mock.Anything).Return(nil)

	h := &TaxHandler{
		TaxService: mockTaxService, // Inject the mocked service
	}

	// Perform the test
	if assert.NoError(t, h.SetPersonalDeduction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestTaxHandler_SetKreceiptDeduction(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 50000}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mocking the TaxService
	mockTaxService := new(MockTaxService)
	mockTaxService.On("SetKReceiptDeduction", mock.Anything).Return(nil)

	h := &TaxHandler{
		TaxService: mockTaxService, // Inject the mocked service
	}

	// Perform the test
	if assert.NoError(t, h.SetKreceiptDeduction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestTaxHandler_TaxCalculationsCSVHandler(t *testing.T) {
	e := echo.New()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxes", "testdata.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("totalIncome,wht,donation\n500000,25000,1000")) // Simulate CSV file content
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mocking the TaxService
	mockTaxService := new(MockTaxService)
	mockTaxService.On("TaxFromFile", mock.Anything).Return([]model.TotalIncomeCsv{{TotalIncome: 500000, WHT: 25000, Donation: 1000}}, nil)

	h := &TaxHandler{
		TaxService: mockTaxService, // Inject the mocked service
	}

	// Perform the test
	if assert.NoError(t, h.TaxCalculationsCSVHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// Optionally verify the response body if needed
	}
}
