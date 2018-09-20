package datastruct

type ErrResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//Websocket resp type: info:1, block:2, tx:3
type WsInfoResp struct {
	Type          int   `json:"type"`
	HighestBolck  int64 `json:"highest_block"`
	LastBlockTime int64 `json:"last_block_time"`
	NodeNum       int   `json:"node_num"`
	TxNum         int64 `json:"tx_num"`
}

type WsBlockResp struct {
	Type        int    `json:"type"`
	Height      int64  `json:"height"`
	ID          string `json:"id"`
	Miner       string `json:"miner"`
	CreatedUsed int64  `json:"created_used"`
	TxCnt       int    `json:"tx_cnt"`
	Time        int64  `json:"time"`
}

type WsTxResp struct {
	Type int `json:"type"`
	DgatewayTx
}

type NodeListResp struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data NodeList `json:"data"`
}
type NodeList struct {
	Count  int64  `json:"count"`
	Online int64  `json:"online"`
	Data   []Node `json:"data"`
}

type BlockListResp struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data BlockList `json:"data"`
}

type BlockList struct {
	Count int64           `json:"count"`
	Data  []DgatewayBlock `json:"data"`
}
type BlockDetailResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data BlockDetail `json:"data"`
}

type BlockDetail struct {
	*DgatewayBlock
	Trans []DgatewayTx `json:"trans"`
}

type CurrentBlockResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Block  `json:"data"`
}

type TxListResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data TxList `json:"data"`
}
type TxList struct {
	Count int64        `json:"count"`
	Data  []DgatewayTx `json:"data"`
}

type TxDetailResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data *DgatewayTx `json:"data"`
}

type InfoTableResp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type InfoOverviewResp struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data WsInfoResp `json:"data"`
}
