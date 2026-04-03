package domain

type AverageResponse struct {
	ReferenceDate string      `json:"referece_date"`
	DigitalCoin   DigitalCoin `json:"digital_coin"`
	AvgType       string      `json:"avg_type"`
	Value         string      `json:"value"`
}
