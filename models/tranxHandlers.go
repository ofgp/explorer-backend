package models

import (
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/dboperation"
)

func GetTxList(page, pageSize int, search string) (*datastruct.TxListResp, error) {
	limit0 := int64((page - 1) * pageSize)
	limit1 := int64(pageSize)
	count, respList, err := dboperation.GetTxListFromMysql(limit0, limit1, search)
	if err != nil {
		return nil, err
	}
	txListResp := &datastruct.TxListResp{
		Code: 0,
		Msg:  "",
		Data: datastruct.TxList{
			Count: count,
			Data:  respList,
		},
	}
	return txListResp, nil
}

func GetTxByHash(dgwHash string) (*datastruct.TxDetailResp, error) {
	resp, err := dboperation.GetTxFromMysql(dgwHash)
	if err != nil {
		return nil, err
	}
	detailResp := &datastruct.TxDetailResp{
		Code: 0,
		Msg:  "",
		Data: resp,
	}
	return detailResp, nil
}
