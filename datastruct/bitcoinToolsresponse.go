package datastruct

//UtxoInfo UTXO数据spend_type: -1:failed; 0:unconfirm; 1: unspent 2:using 3:spent 4:临时占用中
type UtxoInfo struct {
	Address      string `json:"address"`
	VoutTxid     string `json:"vout_txid"`
	VoutIndex    uint32 `json:"vout_index"`
	Value        int64  `json:"value"`
	VoutPkscript string `json:"vout_pkscript"`
	SpendType    int32  `json:"spend_type"`
	VinTxid      string `json:"vin_txid"`
	BlockHeight  int64  `json:"block_height"`
	IsCoinBase   bool   `json:"is_coinbase"`
}

//AddressInfo 地址金额结构体
type AddressInfo struct {
	Address string
	Amount  int64
}

type BatchSendResp struct {
	Code    int            `json:"code"`
	Msg     string         `json:"msg"`
	TxHash  []string       `json:"txhash"`
	Success []*AddressInfo `json:"success"`
	Fail    []*AddressInfo `json:"fail"`
}
