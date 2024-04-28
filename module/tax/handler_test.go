package tax

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestTaxHandler_PostTaxCalculation_Success(t *testing.T) {
	e := echo.New()
	requestBody := `{"totalIncome": 500000, "wht": 25000, "allowances":[{"allowanceType":"k-receipt","amount":50000}]}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedTax := 12345.67
	expectedBrackets := []model.TaxBracket{{Level: "Low", Tax: 1000}}
	mockTaxService := new(MockTaxService)
	mockTaxService.On("CalculateTax", mock.Anything).Return(expectedTax, expectedBrackets, nil)

	h := &TaxHandler{TaxService: mockTaxService}
	if assert.NoError(t, h.PostTaxCalculation(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := `{"tax":12345.7,"taxLevel":[{"level":"Low","tax":1000}]}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestTaxCalculationBindingError(t *testing.T) {
	e := echo.New()
	invalidJSON := `{"totalIncome": "not a float"}`
	req := httptest.NewRequest(http.MethodPost, "/calculateTax", bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &TaxHandler{}

	if err := handler.PostTaxCalculation(c); err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request data")
}

func TestTaxCalculationServiceError(t *testing.T) {
	e := echo.New()
	request := model.TaxRequest{
		TotalIncome: 50000,
		WHT:         5000,
		Allowances: []model.Allowance{
			{
				AllowanceType: "Dummy",
				Amount:        3000,
			},
		},
	}
	reqBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/calculateTax", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	mockTaxService.On("CalculateTax", mock.Anything).Return(0.0, []model.TaxBracket{}, errors.New("calculation error"))

	handler := &TaxHandler{TaxService: mockTaxService}

	if err := handler.PostTaxCalculation(c); err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Tax calculation failed")

	mockTaxService.AssertExpectations(t)
}

func TestTaxHandler_SetPersonalDeduction(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 50000}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	mockTaxService.On("SetPersonalDeduction", mock.Anything).Return(nil)

	h := &TaxHandler{
		TaxService: mockTaxService,
	}

	if assert.NoError(t, h.SetPersonalDeduction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestTaxHandler_SetPersonalDeduction_BindingError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": "invalid"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	h := &TaxHandler{
		TaxService: mockTaxService,
	}

	if assert.NoError(t, h.SetPersonalDeduction(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestTaxHandler_SetPersonalDeduction_InvalidAmount(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 100001}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	h := &TaxHandler{
		TaxService: mockTaxService,
	}

	if assert.NoError(t, h.SetPersonalDeduction(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestTaxHandler_SetKreceiptDeduction(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 50000}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

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

func TestTaxHandler_SetKreceiptDeduction_BindingError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": "invalid"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	h := &TaxHandler{
		TaxService: mockTaxService,
	}

	if assert.NoError(t, h.SetKreceiptDeduction(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestTaxHandler_SetKreceiptDeduction_InvalidAmount(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"amount": 100001}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	h := &TaxHandler{
		TaxService: mockTaxService,
	}

	if assert.NoError(t, h.SetKreceiptDeduction(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
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
	part.Write([]byte("totalIncome,wht,donation\n500000,25000,1000"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTaxService := new(MockTaxService)
	mockTaxService.On("TaxFromFile", mock.Anything).Return([]model.TotalIncomeCsv{{TotalIncome: 500000, WHT: 25000, Donation: 1000}}, nil)

	h := &TaxHandler{
		TaxService: mockTaxService,
	}

	if assert.NoError(t, h.TaxCalculationsCSVHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
