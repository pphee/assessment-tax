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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	tax, err := h.TaxService.CalculateTax(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	res := model.TaxResponse{
		Tax: tax,
	}

	return c.JSON(http.StatusOK, res)
}
