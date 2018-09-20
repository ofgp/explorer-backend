package chainapi

import (
	"github.com/astaxie/beego"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	log "github.com/inconshreveable/log15"
)

//BitCoinClient BTC/BCH RPC操作类
type BitCoinClient struct {
	rpcClient *rpcclient.Client
}

//BlockData 区块数据
type BlockData struct {
	BlockInfo *btcjson.GetBlockVerboseResult
	MsgBolck  *wire.MsgBlock
}

//NewBitCoinClient 创建一个bitcoin操作客户端
func NewBitCoinClient(coinType string) (*BitCoinClient, error) {
	connCfg := &rpcclient.ConnConfig{
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	bc := &BitCoinClient{}

	switch coinType {
	case "btc":
		connCfg.Host = beego.AppConfig.String("btc_rpc_server")
		connCfg.User = beego.AppConfig.String("btc_rpc_user")
		connCfg.Pass = beego.AppConfig.String("btc_rpc_password")
	case "bch":
		connCfg.Host = beego.AppConfig.String("bch_rpc_server")
		connCfg.User = beego.AppConfig.String("bch_rpc_user")
		connCfg.Pass = beego.AppConfig.String("bch_rpc_password")
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Error("GET_BTC_RPC_CLIENT FAIL:", "err", err.Error())
	}

	bc.rpcClient = client

	return bc, err
}

//GetBlockCount 获取当前区块链高度
func (b *BitCoinClient) GetBlockCount() int64 {
	blockHeight, err := b.rpcClient.GetBlockCount()
	if err != nil {
		log.Warn("GET_BLOCK_COUNT FAIL:", "err", err.Error())
		return -1
	}
	return blockHeight
}

//GetRawTransaction 根据txhash从区块链上查询交易数据
func (b *BitCoinClient) GetRawTransaction(txHash string) (*btcutil.Tx, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Warn("NEW_HASH_FAILED:", "err", err.Error(), "hash", txHash)
		return nil, err
	}

	txRaw, err := b.rpcClient.GetRawTransaction(hash)
	if err != nil {
		log.Warn("GetRawTransaction FAILED:", "err", err.Error(), "hash", txHash)
		return nil, err
	}
	return txRaw, nil
}

//GetRawTransactionVerbose 根据txhash从区块链上查询交易数据（包含区块信息）
func (b *BitCoinClient) GetRawTransactionVerbose(txHash string) (*btcjson.TxRawResult, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Warn("NEW_HASH_FAILED:", "err", err.Error(), "hash", txHash)
		return nil, err
	}

	txRaw, err := b.rpcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		log.Warn("GetRawTransactionVerbose FAILED:", "err", err.Error(), "hash", txHash)
		return nil, err
	}
	return txRaw, nil
}

//GetRawMempool 从全节点内存中获取内存中的交易数据
func (b *BitCoinClient) GetRawMempool() ([]*chainhash.Hash, error) {
	result, err := b.rpcClient.GetRawMempool()
	if err != nil {
		log.Warn("GetRawMempool FAILED:", "err", err.Error())
		return nil, err
	}
	return result, err
}

//GetBlockInfoByHeight 根据区块高度获取区块信息
func (b *BitCoinClient) GetBlockInfoByHeight(height int64) *BlockData {
	blockHash, err := b.rpcClient.GetBlockHash(height)
	if err != nil {
		log.Warn("GET_BLOCK_HASH FAIL:", "err", err.Error())
		return nil
	}

	blockVerbose, err := b.rpcClient.GetBlockVerbose(blockHash)
	if err != nil {
		log.Warn("GET_BLOCK_VERBOSE FAIL:", "err", err.Error())
		return nil
	}

	blockEntity, err := b.rpcClient.GetBlock(blockHash)
	if err != nil {
		log.Warn("GET_BLOCK FAIL:", "err", err.Error())
		return nil
	}

	return &BlockData{
		BlockInfo: blockVerbose,
		MsgBolck:  blockEntity,
	}
}

//GetBlockInfoByHash 根据区块hash获取区块信息
func (b *BitCoinClient) GetBlockInfoByHash(hash string) *BlockData {
	blockHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		log.Warn("GET_BLOCK_HASH FAIL:", "err", err.Error())
		return nil
	}

	blockVerbose, err := b.rpcClient.GetBlockVerbose(blockHash)
	if err != nil {
		log.Warn("GET_BLOCK_VERBOSE FAIL:", "err", err.Error(), "hash", blockHash.String())
		return nil
	}

	blockEntity, err := b.rpcClient.GetBlock(blockHash)
	if err != nil {
		log.Warn("GET_BLOCK FAIL:", "err", err.Error())
		return nil
	}

	return &BlockData{
		BlockInfo: blockVerbose,
		MsgBolck:  blockEntity,
	}
}

//SendRawTransaction 发送交易数据到全节点
func (b *BitCoinClient) SendRawTransaction(tx *wire.MsgTx) (*chainhash.Hash, error) {
	return b.rpcClient.SendRawTransaction(tx, true)
}

//EstimateFee 评估交易矿工费
func (b *BitCoinClient) EstimateFee(numBlocks int64) (int64, error) {
	fee, err := b.rpcClient.EstimateFee(numBlocks)
	if err != nil {
		return 0, err
	}

	return int64(fee * 1E8), err

}
