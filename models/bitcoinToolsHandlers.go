package models

import (
	"dgatewayWebBrowser/chainapi"
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/dboperation"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/bchutil"
)

var bt *BitcoinTool

func init() {
	bt = NewBitcoinTool()
}

type BatchSendRequest struct {
	CoinType      string                    `json:"coin_type"`
	PrivateKey    string                    `json:"private_key"`
	ToAddressList []*datastruct.AddressInfo `json:"to"`
	FromAddress   string                    `json:"from"`
}

//BitcoinTool bitcoin工具类
type BitcoinTool struct {
	btcClient *chainapi.BitCoinClient
	bchClient *chainapi.BitCoinClient
	esClient  *dboperation.EsClient
	netParam  *chaincfg.Params
}

//NewBitcoinTool 新建一个bitcoin工具类
func NewBitcoinTool() *BitcoinTool {
	btcClient, err := chainapi.NewBitCoinClient("btc")
	if err != nil {
		return nil
	}

	bchClient, err := chainapi.NewBitCoinClient("bch")
	if err != nil {
		return nil
	}

	esClient := dboperation.NewEsClient()

	var netParam *chaincfg.Params
	switch beego.AppConfig.String("net_param") {
	case "mainnet":
		netParam = &chaincfg.MainNetParams
	case "testnet":
		netParam = &chaincfg.TestNet3Params
	case "regtest":
		netParam = &chaincfg.RegressionNetParams
	}

	bt := &BitcoinTool{
		btcClient: btcClient,
		bchClient: bchClient,
		esClient:  esClient,
		netParam:  netParam,
	}

	return bt
}

func BatchSendBitcoin(request *BatchSendRequest) *datastruct.BatchSendResp {
	result := &datastruct.BatchSendResp{}
	privateKey, err := hex.DecodeString(request.PrivateKey)
	if err != nil {
		result.Code = 2
		result.Msg = "private key err"
		return result
	}

	priKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)

	address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(),
		bt.netParam)

	switch request.CoinType {
	case "btc":
		request.FromAddress = address.EncodeAddress()
	case "bch":
		addr2, _ := bchutil.NewCashAddressPubKeyHash(address.AddressPubKeyHash().ScriptAddress(), bt.netParam)
		request.FromAddress = addr2.EncodeAddress()
	}

	beego.Debug("request.FromAddress", "request.FromAddress", request.FromAddress)

	tempUtxo := bt.esClient.GetUtxoInfoByKey(request.FromAddress, request.CoinType, "address")
	beego.Debug("utxoList", "utxo len", len(tempUtxo))

	utxoMap := make(map[string]*datastruct.UtxoInfo)
	for _, utxo := range tempUtxo {
		if utxo.SpendType > 1 {
			continue
		}

		utxoID := strings.Join([]string{utxo.VoutTxid, strconv.Itoa(int(utxo.VoutIndex))}, "_")
		utxoMap[utxoID] = utxo
	}
	beego.Debug("utxoList", "utxo len", len(utxoMap))

	for offset := 0; offset < len(request.ToAddressList); {
		var tempList []*datastruct.AddressInfo

		if offset+100 >= len(request.ToAddressList) {
			tempList = request.ToAddressList[offset:len(request.ToAddressList)]
		} else {
			tempList = request.ToAddressList[offset : offset+100]
		}

		ret, txhash := bt.SendBitcoin(priKey, request.CoinType, utxoMap, tempList)
		if ret {
			//成功
			result.TxHash = append(result.TxHash, txhash)
			result.Success = append(result.Success, tempList...)
		} else {
			//失败
			result.Code = 1
			result.Msg = "some address failed"
			result.Fail = append(result.Fail, tempList...)
		}
		offset += 100
	}

	return result

}

//DecodeAddress 从地址字符串中decode Address
func (bt *BitcoinTool) DecodeAddress(addr string, coinType string) (btcutil.Address, error) {
	switch coinType {
	case "btc":
		return btcutil.DecodeAddress(addr, bt.netParam)
	case "bch":
		return bchutil.DecodeAddress(addr, bt.netParam)
	default:
		return nil, nil
	}

}

func coinSelect(utxoList []*datastruct.UtxoInfo, targetValue int64) ([]*datastruct.UtxoInfo, int64) {
	var retCoin []*datastruct.UtxoInfo
	var coinValueSum int64

	var lowerCoin []*datastruct.UtxoInfo
	var lowerSum int64

	var coinLowestLarger *datastruct.UtxoInfo
	var coinLowestLargerValue int64 = -1

	for _, utxoInfo := range utxoList {
		if utxoInfo.Value == targetValue {
			retCoin = append(retCoin, utxoInfo)
			coinValueSum += utxoInfo.Value
			return retCoin, coinValueSum
		} else if utxoInfo.Value < targetValue {
			lowerCoin = append(lowerCoin, utxoInfo)
			lowerSum += utxoInfo.Value
		} else if utxoInfo.Value < coinLowestLargerValue || coinLowestLargerValue == -1 {
			coinLowestLarger = utxoInfo
			coinLowestLargerValue = utxoInfo.Value
		}
	}

	if lowerSum == targetValue {
		return lowerCoin, lowerSum
	}

	if lowerSum < targetValue {
		if coinLowestLargerValue == -1 {
			return retCoin, coinValueSum
		}

		retCoin = append(retCoin, coinLowestLarger)
		coinValueSum += coinLowestLarger.Value
		return retCoin, coinValueSum
	}

	nBest := lowerSum
	vBest := make([]int, len(lowerCoin))
	for i := range vBest {
		vBest[i] = 1
	}

	for rep := 0; rep < 100 && nBest != targetValue; rep++ {
		var nTotal int64
		vInclude := make([]int, len(lowerCoin))

		for id, utxoInfo := range lowerCoin {
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(100)%2 == 0 {
				nTotal += utxoInfo.Value
				vInclude[id] = 1

				if nTotal >= targetValue {
					if nTotal <= nBest {
						nBest = nTotal
						copy(vBest, vInclude)
						break
					}
					nTotal -= utxoInfo.Value
					vInclude[id] = 0
				}
			}
		}
	}

	if coinLowestLargerValue != -1 && coinLowestLargerValue-targetValue <= nBest-targetValue {
		retCoin = append(retCoin, coinLowestLarger)
		coinValueSum += coinLowestLarger.Value
		return retCoin, coinValueSum
	}
	for id, utxoInfo := range lowerCoin {
		if vBest[id] == 1 {
			retCoin = append(retCoin, utxoInfo)
		}
	}
	return retCoin, nBest

}

func (bt *BitcoinTool) SendBitcoin(privateKey *btcec.PrivateKey, coinType string, utxoMap map[string]*datastruct.UtxoInfo, toAddressList []*datastruct.AddressInfo) (bool, string) {
	var bitcoinClient *chainapi.BitCoinClient
	switch coinType {
	case "btc":
		bitcoinClient = bt.btcClient
	case "bch":
		bitcoinClient = bt.bchClient
	}

	var utxoList []*datastruct.UtxoInfo
	for _, utxo := range utxoMap {
		if utxo.SpendType > 1 {
			continue
		}
		utxoList = append(utxoList, utxo)
	}
	beego.Debug("utxoList", "utxo len", len(utxoList))

	lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return privateKey, true, nil
	}

	feePerKB, err := bitcoinClient.EstimateFee(1)
	if err != nil {
		beego.Error("get fee/kb failed", "err", err.Error())
		return false, ""
	}
	if feePerKB <= 0 {
		feePerKB = 1000
	}

	var totalValue int64
	for _, addr := range toAddressList {
		totalValue += addr.Amount
	}

	tx := wire.NewMsgTx(2)

	for _, addrInfo := range toAddressList {
		decodeAddr, err := bt.DecodeAddress(addrInfo.Address, coinType)
		if err != nil {
			beego.Warn("DecodeAddress failed", "err", err.Error(), "coinType", coinType)
			return false, ""
		}
		pkScript, err := bchutil.PayToAddrScript(decodeAddr)
		if err != nil {
			beego.Warn("PayToAddrScript failed", "err", err.Error(), "coinType", coinType)
			return false, ""
		}

		vout := wire.TxOut{
			Value:    addrInfo.Amount,
			PkScript: pkScript,
		}
		tx.AddTxOut(&vout)
	}

	minFee := int64(tx.SerializeSize()) * feePerKB / 1000

	for {
		copyTx := tx.Copy()
		selectCoinList, coinSum := coinSelect(utxoList, totalValue+minFee)
		if coinSum < totalValue+minFee {
			return false, ""
		}

		var fromAddress string
		var script []byte
		for _, selectCoin := range selectCoinList {
			fromAddress = selectCoin.Address
			script, _ = hex.DecodeString(selectCoin.VoutPkscript)

			hash, err := chainhash.NewHashFromStr(selectCoin.VoutTxid)
			if err != nil {
				beego.Warn("NEW_HASH_FAILED:", "err", err.Error(), "hash", hash, "coinType", coinType)
				return false, ""
			}

			vin := wire.TxIn{
				PreviousOutPoint: wire.OutPoint{
					Hash:  *hash,
					Index: selectCoin.VoutIndex,
				},
			}

			copyTx.AddTxIn(&vin)
		}

		//评估矿工费
		fee := int64(copyTx.SerializeSize()+110*len(selectCoinList)+40) * feePerKB / 1000
		if fee < 1000 {
			fee = 1000
		}

		beego.Debug("check size and fee", "size", copyTx.SerializeSize(), "fee", fee)

		if coinSum-totalValue > fee {
			hasChange := false
			var smallChange *datastruct.UtxoInfo
			//构造交易
			if coinSum-totalValue-fee > 546 {
				//找零
				decodeAddr, err := bt.DecodeAddress(fromAddress, coinType)
				if err != nil {
					beego.Warn("DecodeAddress failed", "err", err.Error(), "coinType", coinType)
					return false, ""
				}
				pkScript, err := bchutil.PayToAddrScript(decodeAddr)
				if err != nil {
					beego.Warn("PayToAddrScript failed", "err", err.Error(), "coinType", coinType)
					return false, ""
				}

				vout := wire.TxOut{
					Value:    coinSum - totalValue - fee,
					PkScript: pkScript,
				}
				beego.Debug("vout len", "size", vout.SerializeSize())
				tx.AddTxOut(&vout)
				hasChange = true

				smallChange = &datastruct.UtxoInfo{
					Address:      fromAddress,
					VoutIndex:    uint32(len(tx.TxOut) - 1),
					Value:        coinSum - totalValue - fee,
					VoutPkscript: hex.EncodeToString(pkScript),
					SpendType:    0,
				}

			}

			for _, selectCoin := range selectCoinList {
				hash, err := chainhash.NewHashFromStr(selectCoin.VoutTxid)
				if err != nil {
					beego.Warn("NEW_HASH_FAILED:", "err", err.Error(), "hash", hash, "coinType", coinType)
					return false, ""
				}
				selectCoin.SpendType = 2
				beego.Debug("selectCoin", "utxo", selectCoin)

				vin := wire.TxIn{
					PreviousOutPoint: wire.OutPoint{
						Hash:  *hash,
						Index: selectCoin.VoutIndex,
					},
				}

				tx.AddTxIn(&vin)
			}

			//签名
			for i := range tx.TxIn {
				if coinType == "bch" {
					sigScript, _ := bchutil.SignTxOutput(bt.netParam,
						tx, i, script, txscript.SigHashAll,
						txscript.KeyClosure(lookupKey), nil, nil, selectCoinList[i].Value)

					tx.TxIn[i].SignatureScript = sigScript
					beego.Debug("sign len", "len", len(sigScript))
				} else {
					sigScript, _ := txscript.SignTxOutput(bt.netParam,
						tx, i, script, txscript.SigHashAll,
						txscript.KeyClosure(lookupKey), nil, nil)
					tx.TxIn[i].SignatureScript = sigScript
				}
			}

			beego.Debug("check size and fee", "size", tx.SerializeSize(), "fee", fee)
			//发送交易
			txHash, err := bitcoinClient.SendRawTransaction(tx)
			if err != nil {
				beego.Warn("send err", "err", err.Error())
				return false, ""
			}

			if hasChange {
				smallChange.VoutTxid = txHash.String()
				utxoID := strings.Join([]string{smallChange.VoutTxid, strconv.Itoa(int(smallChange.VoutIndex))}, "_")
				utxoMap[utxoID] = smallChange
			}

			return true, txHash.String()

		}
		minFee = fee
	}

}
