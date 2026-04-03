package domain

type AverageResponse struct {
	ReferenceDate string      `json:"referece_date"`
	DigitalCoin   DigitalCoin `json:"digital_coin"`
	AvgType       string      `json:"avg_type"`
	Days          int         `json:"days"`
	Value         string      `json:"value"`
}
