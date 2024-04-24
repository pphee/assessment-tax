package model

import (
	"encoding/json"
	"fmt"
)

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxResponse struct {
	Tax float64 `json:"-"`
}

func (tr TaxResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Tax json.Number `json:"tax"`
	}{
		Tax: json.Number(fmt.Sprintf("%.1f", tr.Tax)),
	})
}
