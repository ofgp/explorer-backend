package dboperation

import (
	"context"
	"dgatewayWebBrowser/datastruct"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/olivere/elastic"
)

type BtcTxInfo struct {
	BlockHash   string      `json:"block_hash"`
	BlockHeight int64       `json:"block_height"`
	BlockTime   int64       `json:"block_time"`
	Txid        string      `json:"txid"`
	TxidIndex   int32       `json:"txid_index"`
	Version     int32       `json:"version"`
	Size        int32       `json:"size"`
	Locktime    uint32      `json:"locktime"`
	Status      int32       `json:"status"`
	Vin         []*VinInfo  `json:"vin"`
	Vout        []*VoutInfo `json:"vout"`
}

//VinInfo 交易输入数据
type VinInfo struct {
	PreTxid         string `json:"pre_txid"`
	PreTxidVout     uint32 `json:"pre_txid_vout"`
	PreVoutValue    int64  `json:"pre_vout_value"`
	SignatureScript string `json:"signature_script"`
	InputAddress    string `json:"input_address"`
	Sequence        uint32 `json:"sequence"`
}

//VoutInfo 交易输出数据
type VoutInfo struct {
	VoutValue     int64  `json:"vout_value"`
	OutputAddress string `json:"output_address"`
	PkScript      string `json:"pk_script"`
}

//EOS tx
type EosTx struct {
	BlockNum       uint32 `json:"block_num"`
	BlockID        string `json:"block_id"`
	BlockTimeStamp int64  `json:"block_timestamp"`
	Producer       string `json:"producer"`
	TxID           string `json:"tx_id"`
	Status         string `json:"status"`
	//CpuUsageUs uint32
	//NetUsageWords eos.Varuint32
	Expiration     int64    `json:"expiration"`
	RefBlockNum    uint16   `json:"ref_block_num"`
	RefBlockPrefix uint32   `json:"ref_block_prefix"`
	DelaySec       uint32   `json:"delay_sec"`
	Actions        []Action `json:"actions"`
}

// EOS tx Action
type Action struct {
	Account string          `json:"account"`
	Name    string          `json:"name"`
	Data    json.RawMessage `json:"data"`
}

// destorytoken method params
type DestoryTokenParams struct {
	User   string `json:"user"`
	Amount uint32 `json:"amount"`
	Memo   string `json:"memo"`
}

type EsClient struct {
	sync.Mutex
	client *elastic.Client
	ctx    context.Context
}

func NewEsClient() *EsClient {
	ctx := context.Background()
	esClient, err := elastic.NewClient(
		elastic.SetURL(beego.AppConfig.String("elasticsearchurl")),
		elastic.SetSniff(false),
	)
	if err != nil {
		beego.Error(err)
		panic("Wrong ESCLient ")
	}
	return &EsClient{
		ctx:    ctx,
		client: esClient,
	}
}

func (ec *EsClient) getBTCTxByHash(hash, chainCode string) (out *BtcTxInfo, err error) {
	idsQuery := elastic.NewMatchQuery("txid", hash)
	index := []string{}
	if chainCode == "btc" {
		index = beego.AppConfig.Strings("elasticsearchbtcindex")
	} else if chainCode == "bch" {
		index = beego.AppConfig.Strings("elasticsearchbchindex")
	} else {
		beego.Error("Invalid chainCode", chainCode)
		return nil, fmt.Errorf("invalid chainCode")
	}
	searchResults, err := ec.client.Search().
		Index(index...).
		Query(idsQuery).
		Pretty(true).
		Do(ec.ctx)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if searchResults.Hits.TotalHits > 0 {
		var tx BtcTxInfo
		err := json.Unmarshal(*searchResults.Hits.Hits[0].Source, &tx)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		return &tx, nil
	}
	return nil, fmt.Errorf("No Results:", hash)
}

//get btc tx fee, search via elasticsearch
func (ec *EsClient) GetBtcTxFee(hash, chainCode string) (fee int64, err error) {
	btcTx, err := ec.getBTCTxByHash(hash, chainCode)
	if err != nil {
		beego.Error(err)
		return 0, nil
	}

	input := int64(0)
	output := int64(0)
	for _, in := range btcTx.Vin {
		//coinbase
		if in.PreTxid == "0000000000000000000000000000000000000000000000000000000000000000" {
			//log.Info
			return 0, nil
		}
		input = input + in.PreVoutValue
	}
	for _, out := range btcTx.Vout {
		output = output + out.VoutValue
	}
	return input - output, nil
}

//get btc tx From, search via elasticsearch
func (ec *EsClient) GetBtcTxFrom(hash, chainCode string) (fromAddr string, err error) {
	btcTx, err := ec.getBTCTxByHash(hash, chainCode)
	if err != nil {
		beego.Error(err)
		return "", nil
	}
	addrlist := []string{}
	for _, in := range btcTx.Vin {
		//coinbase
		if in.PreTxid == "0000000000000000000000000000000000000000000000000000000000000000" {
			//log.Info
			return "", nil
		}
		addrlist = append(addrlist, in.InputAddress)
	}
	return strings.Join(addrlist, ","), nil
}

//GetUtxoInfoByKey 查询某个交易关联的utxo
func (ec *EsClient) GetUtxoInfoByKey(value string, coinType string, key string) []*datastruct.UtxoInfo {
	var utxoList []*datastruct.UtxoInfo

	indexName := strings.Join([]string{coinType, "utxo"}, "_")

	query := elastic.NewMatchQuery(key, value)
	scroller := ec.client.Scroll().Index(indexName).
		Type("_doc").
		Query(query)

	for {
		res, err := scroller.Do(ec.ctx)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF || res == nil {
			break
		}
		var tempUtxo datastruct.UtxoInfo
		for _, item := range res.Each(reflect.TypeOf(tempUtxo)) {
			if t, ok := item.(datastruct.UtxoInfo); ok {
				utxoList = append(utxoList, &t)
			}
		}
	}
	scroller.Clear(ec.ctx)

	return utxoList
}

//GetEosDestoryTokenTxFrom query certain destoryTokenTx of EOS
func (ec *EsClient) GetEosDestoryTokenTxFrom(txHash string) (string, error) {
	eosTx, err := ec.getEosXinTxByHash(txHash)
	if err != nil {
		return "", err
	}
	addrList := []string{}
	for _, action := range eosTx.Actions {
		if action.Name == "destroytoken" {
			destoryTokenParams := &DestoryTokenParams{}
			err := json.Unmarshal(action.Data, destoryTokenParams)
			if err != nil {
				beego.Error(err)
				continue
			}
			addrList = append(addrList, destoryTokenParams.User)
		}
	}
	return strings.Join(addrList, ","), nil
}

func (ec *EsClient) getEosXinTxByHash(txHash string) (*EosTx, error) {
	idsQuery := elastic.NewMatchQuery("tx_id", txHash)
	index := beego.AppConfig.String("elasticsearch_eos_xin_index")
	searchResults, err := ec.client.Search().
		Index(index).
		Query(idsQuery).
		Pretty(true).
		Do(ec.ctx)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if searchResults.Hits.TotalHits > 0 {
		var out EosTx
		err := json.Unmarshal(*searchResults.Hits.Hits[0].Source, &out)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		return &out, nil
	}
	return nil, fmt.Errorf("No Results:", txHash)
}
