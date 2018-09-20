package datastruct

import (
	"time"
)

type DgatewayBlock struct {
	Id          int64  `json:"-"`
	Height      int64  `json:"height"`
	Hash        string `json:"hash"`
	PreId       string `json:"pre_id"`
	TxCnt       int    `json:"tx_cnt"`
	Time        int64  `json:"time"`
	Size        int64  `json:"size"`
	CreatedUsed int64  `json:"created_used"`
	Miner       string `json:"miner"`
}

type DgatewayTx struct {
	Id              int64  `json:"-"`
	FromTxHash      string `json:"from_tx_hash"`
	ToTxHash        string `json:"to_tx_hash"`
	DgwHash         string `json:"dgw_hash"`
	FromChain       string `json:"from_chain"`
	ToChain         string `json:"to_chain"`
	Time            int64  `json:"time"`
	Block           string `json:"block"`
	BlockHeight     int64  `json:"block_height"`
	Amount          int64  `json:"amount"`
	FromAddrs       string `orm:"size:(1000)" json:"from_addrs"`
	ToAddrs         string `orm:"size(1000)" json:"to_addrs"`
	FromFee         int64  `json:"from_fee"`
	DgwFee          int64  `json:"dgw_fee"`
	ToFee           int64  `json:"to_fee"`
	TokenCode       uint32 `json:"token_code"`
	AppCode         uint32 `json:"app_code"`
	TokenSymbol     string `json:"token_symbol"`
	TokenDecimals   int    `json:"token_decimals"`
	FinalAmount     int64  `json:"final_amount"`
	ToTokenSymbol   string `json:"to_token_symbol"`
	ToTokenDecimals int    `json:"to_token_decimals"`
}

//insert into dgateway_tx_statistics (time, amount, count, from_chain) values ("2018-07-08", 2000, 3, "btc");
type DgatewayTxStatistics struct {
	Id             int64     `json:"-"`
	Time           time.Time `json:"time"`
	Amount         int64     `json:"amount"`
	CurrencyAmount float64   `json:"curency_amount"`
	Count          int64     `json:"count"`
	Symbol         string    `json:"symbol"`
}
type DgatewayScanHeight struct {
	Id     int64  `json:"-"`
	Name   string `json:"name"`
	Height int64  `json:"height"`
}
type DgatewayTokenInfo struct {
	Id              int64  `json:"-"`
	Chain           string `json:"chain"`
	TokenCode       uint32 `json:"token_code"`
	Symbol          string `json:"symbol"`
	Decimals        int    `json:"decimals"`
	RelateChain     string `json:"relate_chain"`      //关联chain，eth合约token使用
	RelateTokenCode uint32 `json:"relate_token_code"` //关联token， eth合约token关联的其他链上的token，价值对等。
}
