package tax

import (
	"github.com/labstack/echo/v4"
	"github.com/pphee/assessment-tax/internal/model"
	"github.com/pphee/assessment-tax/utils"
	"mime/multipart"
	"net/http"
)

type TaxHandler struct {
	TaxService *TaxService
}

func NewTaxHandler(service *TaxService) *TaxHandler {
	return &TaxHandler{TaxService: service}
}

func (h *TaxHandler) PostTaxCalculation(c echo.Context) error {
	var req model.TaxRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request data: " + err.Error()})
	}

	tax, taxBrackets, err := h.TaxService.CalculateTax(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Tax calculation failed: " + err.Error()})
	}

	res := model.TaxResponse{
		Tax:       tax,
		TaxLevels: taxBrackets,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TaxHandler) SetPersonalDeduction(c echo.Context) error {
	var req model.AdminRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON format or data types"})
	}

	if req.Amount < 10000 || req.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Amount must be between 10,000 and 100,000"})
	}

	if err := h.TaxService.SetPersonalDeduction(req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to set personal deduction: " + err.Error()})
	}

	res := model.AdminPersonalDeductionResponse{
		Amount: req.Amount,
	}

	return c.JSON(http.StatusOK, res)
}

// SetKreceiptDeduction
func (h *TaxHandler) SetKreceiptDeduction(c echo.Context) error {
	var req model.AdminRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON format or data types"})
	}

	if req.Amount < 0 || req.Amount > 100000 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Amount must be between 0 and 100,000"})
	}

	if err := h.TaxService.SetKReceiptDeduction(req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to set K-receipt deduction: " + err.Error()})
	}

	res := model.AdminKReceiptDeductionResponse{
		Amount: req.Amount,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TaxHandler) TaxCalculationsCSVHandler(c echo.Context) error {
	file, err := c.FormFile("taxes")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No file uploaded"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error opening file"})
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	records, err := h.TaxService.TaxFromFile(src)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error reading file"})
	}

	var taxDetails []model.TaxDetail
	for _, record := range records {
		PersonalDeductionDefault := 60000.00
		taxableIncome := record.TotalIncome - PersonalDeductionDefault - record.Donation
		tax, _ := utils.CalculateIncomeTaxDetailed(taxableIncome)
		netTax := tax - record.WHT

		taxRefund := 0.0
		if netTax < 0 {
			taxRefund = -netTax
			netTax = 0
		}

		taxDetails = append(taxDetails, model.TaxDetail{
			TotalIncome: record.TotalIncome,
			Tax:         netTax,
			TaxRefund:   taxRefund,
		})
	}

	response := model.TaxResponseCSV{
		Taxes: taxDetails,
	}

	return c.JSON(http.StatusOK, response)
}
