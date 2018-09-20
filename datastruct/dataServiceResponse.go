package datastruct

type TokenCurrencyResp struct {
	Price  float64 `json:"price"`
	Target string  `json:"target"`
	Unit   string  `json:"unit"`
	Time   int64   `json:"time"`
}
