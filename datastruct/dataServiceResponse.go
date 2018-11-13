package datastruct

/*deprecated
type TokenCurrencyResp struct {
	Price  float64 `json:"price"`
	Target string  `json:"target"`
	Unit   string  `json:"unit"`
	Time   int64   `json:"time"`
}
*/

type DataCommonResp struct {
	Success bool   `json:"success"`
	SData   string `json:"sData"`
	Code    int64  `json:"code"`
	Msg     string `json:"msg"`
}

type TokenCurrencyResp struct {
	DataCommonResp
	Data PriceData `json:"data"`
}

type PriceData struct {
	TokenName         string  `json:"tokenName"`
	Increase          string  `json:"increase"`
	CurrentTokenPrice float64 `json:"currentTokenPrice"`
}

type ExchangeRateResp struct {
	DataCommonResp
	Data []ExchangeRateData `json:"data"`
}

type ExchangeRateData struct {
	Name   string  `json:"name"`
	CName  string  `json:"cName"` //中文名称
	Symbol string  `json:"symbol"`
	Rate   float64 `json:"rate"`
}

type DataCryptoKeyResponse struct {
	DataCommonResp
	Data string `json:"data"`
}
