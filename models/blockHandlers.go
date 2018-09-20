package models

import (
	"dgatewayWebBrowser/chainapi"
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/dboperation"
)

func GetBlockList(page, pageSize int, search string) (*datastruct.BlockListResp, error) {
	//分页
	limit0 := int64((page - 1) * pageSize)
	limit1 := int64(pageSize)
	count, respList, err := dboperation.GetBlockListFromMysql(limit0, limit1, search)
	if err != nil {
		return nil, err
	}
	bListResp := &datastruct.BlockListResp{
		Code: 0,
		Msg:  "",
		Data: datastruct.BlockList{
			Count: count,
			Data:  respList,
		},
	}
	return bListResp, nil
}

func GetBlockDetail(height string) (*datastruct.BlockDetailResp, error) {
	blockData, err := dboperation.GetBlockDataFromMysql(height)
	if err != nil {
		return nil, err
	}
	txData, err := dboperation.GetBlockTxFromMysql(height)
	if err != nil {
		return nil, err
	}
	data := datastruct.BlockDetail{
		DgatewayBlock: blockData,
		Trans:         txData,
	}
	resp := &datastruct.BlockDetailResp{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	return resp, nil
}

func GetCurrentBlock() (*datastruct.CurrentBlockResp, error) {
	blockData, err := chainapi.GetCurrentBlock()
	if err != nil {
		return nil, err
	}
	data := &datastruct.CurrentBlockResp{
		Code: 0,
		Msg:  "",
		Data: *blockData.Data,
	}
	return data, nil
}
