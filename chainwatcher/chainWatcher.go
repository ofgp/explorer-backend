/*
观察链上新块的产生
1.向浏览器推送新块及新块的交易
2.将观察到的交易写到mysql
*/
package chainwatcher

import (
	"dgatewayWebBrowser/chainapi"
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/dboperation"
	"dgatewayWebBrowser/models"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

func StartWatch() {
	scanHeight, err := beego.AppConfig.Int64("chainScanHeight")
	if err != nil {
		panic(err)
	}
	chainWatcher := &ChainWatcher{
		ScanHeight:   scanHeight,
		BlockChannel: make(chan *datastruct.SingleBlockResp, 100),
		esClient:     dboperation.NewEsClient(),
		redisPool:    dboperation.NewRedisPool(),
	}
	chainWatcher.WatchBlock()
	chainWatcher.StoreTx()
}

type ChainWatcher struct {
	HighestBlock int64
	ScanHeight   int64
	BlockChannel chan *datastruct.SingleBlockResp
	esClient     *dboperation.EsClient
	redisPool    *redis.Pool
}

//WatcheBlock  开始观察区块信息，并进行数据库读写操作
func (cw *ChainWatcher) WatchBlock() {
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer func() {
			wg.Done()
		}()
		//如果初始观察高度是0，从mysql中获取持久化高度
		if cw.ScanHeight == 0 {
			n, err := dboperation.ScanHeightRead()
			if err != nil {
				panic("Get ScanHeight Err")
			}
			if n == 0 {
				cw.ScanHeight = 0
			} else {
				cw.ScanHeight = n
			}
		}
		for {
			blockHeight := chainapi.GetHighestBlockHeight()
			beego.Info("highestBlock:", blockHeight)
			cw.HighestBlock = blockHeight
			if blockHeight <= cw.ScanHeight {
				time.Sleep(time.Duration(10) * time.Second)
				continue
			}
			for {
				blockData, err := chainapi.GetBlockByHeight(cw.ScanHeight)

				if err != nil {
					break
				}
				cw.BlockChannel <- blockData
				cw.ScanHeight++
				//每100个块做一次高度持久化存储
				if cw.ScanHeight%100 == 0 {
					dboperation.ScanHeightUpdate(cw.ScanHeight)
				}
				if cw.ScanHeight > cw.HighestBlock {
					time.Sleep(time.Duration(10) * time.Second)
					break
				}
			}

		}
	}()
	wg.Wait()
}

func (cw *ChainWatcher) StoreTx() {
	blockChan := cw.BlockChannel
	for {
		select {
		case blockData := <-blockChan:
			cw.processBlockData(blockData)
		}
	}
}

func (cw *ChainWatcher) processBlockData(blockData *datastruct.SingleBlockResp) bool {
	//用于浏览器首页block和info
	models.BroadCastBlockChan <- blockData.Data

	//存储block信息到mysql
	block := &datastruct.DgatewayBlock{
		Height:      blockData.Data.Height,
		Hash:        blockData.Data.ID,
		PreId:       blockData.Data.PreID,
		TxCnt:       blockData.Data.TxCnt,
		Time:        blockData.Data.Time,
		Size:        blockData.Data.Size,
		CreatedUsed: blockData.Data.CreatedUsed,
		Miner:       blockData.Data.Miner,
	}
	beego.Info("current watch block height:", block.Height)
	if !dboperation.StoreBlockToMysql(block) {
		panic("Store block data to mysql Err")
	}
	for _, tx := range blockData.Data.TXS {

		//获取"转入链"交易费
		toFee := int64(0)
		if tx.To == "btc" || tx.To == "bch" {
			toFee, _ = cw.esClient.GetBtcTxFee(tx.ToTxHash, tx.To)

		} else if tx.To == "eos" {
			toFee = 0
		} else {
			toFee, _ = chainapi.GetEthMinerFee(tx.ToTxHash)
		}
		//获取"转出链"交易费
		fromFee := int64(0)
		if tx.From == "btc" || tx.From == "bch" {
			fromFee, _ = cw.esClient.GetBtcTxFee(tx.FromTxHash, tx.From)

		} else if tx.From == "eos" {
			fromFee = 0
		} else {
			fromFee, _ = chainapi.GetEthMinerFee(tx.FromTxHash)
		}

		//获取"转入链FromAddr"
		fromAddrString := ""
		if tx.From == "btc" || tx.From == "bch" {
			fromAddrString, _ = cw.esClient.GetBtcTxFrom(tx.FromTxHash, tx.From)
		} else if tx.From == "eos" {
			// xinplayer的熔币交易，方法名为destorytoken
			fromAddrString, _ = cw.esClient.GetEosDestoryTokenTxFrom(tx.FromTxHash)
		} else {
			fromAddrString, _ = chainapi.GetEthTxFrom(tx.FromTxHash)
		}
		toAddrString := strings.Join(tx.ToAddrs, ",")

		//获取amount对应的代币信息, 配置在mysql(dgateway_token_info), eth还需要动态更新
		tokenSymbol, tokenDecimals := getTokenSymbolAndDecimals(tx, "from")
		toTokenSymbol, toTokenDecimals := getTokenSymbolAndDecimals(tx, "to")

		newTx := &datastruct.DgatewayTx{
			FromTxHash:      tx.FromTxHash,
			ToTxHash:        tx.ToTxHash,
			DgwHash:         tx.DGWHash,
			FromChain:       tx.From,
			ToChain:         tx.To,
			Time:            tx.Time,
			Block:           tx.Block,
			BlockHeight:     tx.BlockHeight,
			Amount:          tx.Amount,
			FromAddrs:       fromAddrString,
			ToAddrs:         toAddrString,
			FromFee:         fromFee,
			DgwFee:          tx.DGWFee,
			ToFee:           toFee,
			TokenCode:       tx.TokenCode,
			AppCode:         tx.AppCode,
			TokenSymbol:     tokenSymbol,
			TokenDecimals:   tokenDecimals,
			FinalAmount:     tx.FinalAmount,
			ToTokenSymbol:   toTokenSymbol,
			ToTokenDecimals: toTokenDecimals,
		}
		models.BroadCastTxChan <- newTx
		if !dboperation.StoreTxToMysql(newTx) {
			panic("Store Tx To Mysql Err")
		}
		//推送消息，用于钱包服务
		message, _ := json.Marshal(newTx)
		cw.redisPushTx(string(message))
	}
	return true
}

//根据交易信息获取amount对应的token symbol 和decimals,
//如果eth的token信息获取不到，则从链上抓取，并且更新到mysql
//bch对应主链:token_code, eos,eth对应侧链: app_code
func getTokenSymbolAndDecimals(tx datastruct.Transaction, chain string) (string, int) {
	if chain == "from" {
		if tx.From == "eth" {
			tokenInfo, err := dboperation.GetTokenInfo(tx.From, tx.AppCode)
			if err != nil {
				symbol, decimals := chainapi.GetDgateWayAppSymbolAndDecimal(tx.AppCode)
				dboperation.AddTokenInfo("eth", symbol, tx.To, decimals, tx.AppCode, tx.TokenCode)
				return symbol, decimals
			} else {
				return tokenInfo.Symbol, tokenInfo.Decimals
			}
		} else if tx.From == "eos" {
			tokenInfo, err := dboperation.GetTokenInfo(tx.From, tx.AppCode)
			if err != nil {
				panic("get token info failed")
			}
			return tokenInfo.Symbol, tokenInfo.Decimals
		} else {
			tokenInfo, err := dboperation.GetTokenInfo(tx.From, tx.TokenCode)
			if err != nil {
				panic("get token info failed")
			}
			return tokenInfo.Symbol, tokenInfo.Decimals
		}
	} else {
		if tx.To == "eth" {
			tokenInfo, err := dboperation.GetTokenInfo(tx.To, tx.AppCode)
			if err != nil {
				symbol, decimals := chainapi.GetDgateWayAppSymbolAndDecimal(tx.AppCode)
				dboperation.AddTokenInfo("eth", symbol, tx.From, decimals, tx.AppCode, tx.TokenCode)
				return symbol, decimals
			} else {
				return tokenInfo.Symbol, tokenInfo.Decimals
			}
		} else if tx.To == "eos" {
			tokenInfo, err := dboperation.GetTokenInfo(tx.To, tx.AppCode)
			if err != nil {
				panic("get token info failed")
			}
			return tokenInfo.Symbol, tokenInfo.Decimals
		} else {
			tokenInfo, err := dboperation.GetTokenInfo(tx.To, tx.TokenCode)
			if err != nil {
				panic("get token info failed")
			}
			return tokenInfo.Symbol, tokenInfo.Decimals
		}
	}
}

func (cw *ChainWatcher) redisPushTx(message string) bool {
	conn := cw.redisPool.Get()
	defer conn.Close()

	res, err := conn.Do("LPUSH", beego.AppConfig.String("redisPushKey"), message)
	beego.Info("push redis:", message, "result:", res)
	if err != nil {
		beego.Error("push redis message failed", "msg", message)

		return false
	}
	return true
}
