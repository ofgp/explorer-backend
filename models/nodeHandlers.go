package models

import (
	"dgatewayWebBrowser/chainapi"
	"dgatewayWebBrowser/datastruct"
)

var mockResp = &datastruct.NodeListResp{
	Code: 0,
	Msg:  "",
	Data: datastruct.NodeList{
		Count:  2,
		Online: 1,
		Data: []datastruct.Node{
			datastruct.Node{
				IP:        "127.0.0.1:9999",
				HostName:  "yyy",
				IsLeader:  true,
				IsOnline:  true,
				FiredCnt:  0,
				EthHeight: 111,
				BchHeight: 999,
				BtcHeight: 777,
			},
			datastruct.Node{
				IP:        "127.0.0.1:8888",
				HostName:  "yyyong",
				IsLeader:  true,
				IsOnline:  false,
				FiredCnt:  0,
				EthHeight: 112,
				BchHeight: 998,
				BtcHeight: 779,
			},
		},
	},
}

func GetNodeList() (*datastruct.NodeListResp, error) {
	resp, err := chainapi.GetNodes()
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return &datastruct.NodeListResp{
			Code: -1,
			Msg:  resp.Msg,
		}, nil
	}

	onlineNode := int64(0)
	for _, node := range resp.Data {
		if node.IsOnline == true {
			onlineNode++
		}
	}
	return &datastruct.NodeListResp{
		Code: 0,
		Msg:  "",
		Data: datastruct.NodeList{
			Count:  int64(len(resp.Data)),
			Online: onlineNode,
			Data:   resp.Data,
		},
	}, nil
}
