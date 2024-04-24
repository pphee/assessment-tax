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
	Tax       float64      `json:"-"`
	TaxLevels []TaxBracket `json:"taxLevel"`
}

func (tr TaxResponse) MarshalJSON() ([]byte, error) {
	taxLevels := make([]struct {
		Level string      `json:"level"`
		Tax   json.Number `json:"tax"`
	}, len(tr.TaxLevels))

	for i, tl := range tr.TaxLevels {
		taxLevels[i] = struct {
			Level string      `json:"level"`
			Tax   json.Number `json:"tax"`
		}{
			Level: tl.Level,
			Tax:   json.Number(fmt.Sprintf("%.1f", tl.Tax)),
		}
	}

	return json.Marshal(&struct {
		Tax       json.Number `json:"tax"`
		TaxLevels []struct {
			Level string      `json:"level"`
			Tax   json.Number `json:"tax"`
		} `json:"taxLevel"`
	}{
		Tax:       json.Number(fmt.Sprintf("%.1f", tr.Tax)),
		TaxLevels: taxLevels,
	})
}

type TaxBracket struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type AdminRequest struct {
	Amount float64 `json:"amount"`
}

type AdminPersonalDeductionResponse struct {
	Amount float64 `json:"-"`
}

func (ar AdminPersonalDeductionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Amount json.Number `json:"personalDeduction"`
	}{
		Amount: json.Number(fmt.Sprintf("%.1f", ar.Amount)),
	})
}

type AdminKReceiptDeductionResponse struct {
	Amount float64 `json:"-"`
}

func (ar AdminKReceiptDeductionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Amount json.Number `json:"kReceipt"`
	}{
		Amount: json.Number(fmt.Sprintf("%.1f", ar.Amount)),
	})

}

type TotalIncomeCsv struct {
	TotalIncome float64 `csv:"totalIncome"`
	WHT         float64 `csv:"wht"`
	Donation    float64 `csv:"donation"`
}

type TaxDetail struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}

type TaxResponseCSV struct {
	Taxes []TaxDetail `json:"taxes"`
}
