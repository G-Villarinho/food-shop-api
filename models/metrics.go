package models

type MonthlyMetricsResponse struct {
	Amount            float64 `json:"amount"`
	DiffFromLastMonth float64 `json:"diffFromLastMonth"`
}
