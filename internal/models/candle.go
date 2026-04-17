package models

type Candle struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type MercadoBitcoinResponse struct {
	Status string   `json:"s,omitempty"`
	Close  []string `json:"c"`
	High   []string `json:"h"`
	Low    []string `json:"l"`
	Open   []string `json:"o"`
	Time   []int64  `json:"t"`
	Volume []string `json:"v"`
}
