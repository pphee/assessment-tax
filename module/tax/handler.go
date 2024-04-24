package tax

import (
	"github.com/labstack/echo/v4"
	"github.com/pphee/assessment-tax/internal/model"
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

	return c.JSON(http.StatusOK, echo.Map{"personalDeduction": req.Amount})
}
