package chainapi

import (
	"dgatewayWebBrowser/datastruct"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
)

//获取当前最新区块信息
func GetCurrentBlock() (*datastruct.SingleBlockResp, error) {
	data, err := httpGet(beego.AppConfig.String("chainurl")+"/block/current", map[string]string{})
	if err != nil {
		return nil, err
	}
	blockResp := datastruct.SingleBlockResp{}
	errM := json.Unmarshal(data, &blockResp)

	if errM != nil {
		beego.Error(errM)
		return nil, err
	}

	return &blockResp, nil
}

func GetHighestBlockHeight() (height int64) {
	blockResp, err := GetCurrentBlock()
	if err != nil {
		return 0
	}
	return blockResp.Data.Height
}

//获取区块区间信息 >=start, <end
func GetBlocks(start, end int64) (*[]datastruct.BulkBlockResp, error) {
	startstring := strconv.FormatInt(start, 10)
	endstring := strconv.FormatInt(end, 10)
	params := map[string]string{"start": startstring, "end": endstring}
	data, err := httpGet(beego.AppConfig.String("chainurl")+"/blocks", params)
	if err != nil {
		return nil, err
	}
	blockResp := []datastruct.BulkBlockResp{}
	errM := json.Unmarshal(data, &blockResp)
	if errM != nil {
		return nil, err
	}

	return &blockResp, nil
}

//获取指定区块信息
func GetBlockByHeight(height int64) (*datastruct.SingleBlockResp, error) {
	heightString := strconv.FormatInt(height, 10)
	data, err := httpGet(beego.AppConfig.String("chainurl")+"/block/height/"+heightString, map[string]string{})
	if err != nil {
		return nil, err
	}
	blockResp := datastruct.SingleBlockResp{}
	errM := json.Unmarshal(data, &blockResp)
	if errM != nil {
		return nil, err
	}
	return &blockResp, nil
}

//获取当前全部Node信息
func GetNodes() (*datastruct.BulkNodeResp, error) {
	data, err := httpGet(beego.AppConfig.String("chainurl")+"/nodes", map[string]string{})
	if err != nil {
		return nil, err
	}
	nodeResp := datastruct.BulkNodeResp{}
	errM := json.Unmarshal(data, &nodeResp)
	if errM != nil {
		return nil, err
	}
	return &nodeResp, err
}
func GetNodeNum() (num int, err error) {
	nodeResp, err := GetNodes()
	if err != nil {
		return 0, err
	}
	return len(nodeResp.Data), nil
}

//获取指定交易详情
func GetTransaction(TxID string) (*datastruct.SingleTxResp, error) {
	data, err := httpGet(beego.AppConfig.String("chainurl")+"/transaction/"+TxID, map[string]string{})
	if err != nil {
		return nil, err
	}

	tranxResp := datastruct.SingleTxResp{}
	errM := json.Unmarshal(data, &tranxResp)
	if errM != nil {
		return nil, err
	}
	return &tranxResp, nil
}
