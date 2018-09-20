package datastruct

type Block struct {
	Height      int64         `json:"height"`
	ID          string        `json:"id"`
	PreID       string        `json:"pre_id"`
	TxCnt       int           `json:"tx_cnt"`
	TXS         []Transaction `json:"txs"`
	Time        int64         `json:"time"`
	Size        int64         `json:"size"`
	CreatedUsed int64         `json:"created_used"`
	Miner       string        `json:"miner"`
}

type Transaction struct {
	FromTxHash  string   `json:"from_tx_hash"`
	ToTxHash    string   `json:"to_tx_hash"`
	DGWHash     string   `json:"dgw_hash"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Time        int64    `json:"time"`
	Block       string   `json:"block"`
	BlockHeight int64    `json:"block_height"`
	Amount      int64    `json:"amount"`
	ToAddrs     []string `json:"to_addrs"`
	FromFee     int64    `json:"from_fee"`
	DGWFee      int64    `json:"dgw_fee"`
	ToFee       int64    `json:"to_fee"`
	TokenCode   uint32   `json:"token_code"`
	AppCode     uint32   `json:"app_code"`
	FinalAmount int64    `json:"final_amount"`
}

type Node struct {
	IP        string `json:"ip"`
	HostName  string `json:"host_name"`
	IsLeader  bool   `json:"is_leader"`
	IsOnline  bool   `json:"is_online"`
	FiredCnt  int32  `json:"fired_cnt"`
	EthHeight int64  `json:"eth_height"`
	BchHeight int64  `json:"bch_height"`
	BtcHeight int64  `json:"btc_height"`
}

type SingleBlockResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *Block `json:"data"`
}

type BulkBlockResp struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data []Block `json:"data"`
}

type SingleTxResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data Transaction `json:"data"`
}

type BulkNodeResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Node `json:"data"`
}
