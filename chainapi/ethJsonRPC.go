package chainapi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ybbus/jsonrpc"
)

//EthTxInfo :jsonrpc EthTx Response
type EthTxInfo struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNum         string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	Gas              string `json:"gas"`
	Input            string `json:"string"`
}

//EthTxReceipt :jsonrpc EthTxReceipt Response
type EthTxReceipt struct {
	BlockHash         string `json:"blockHash"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	ContractAddr      string `json:"contractAddress"`
	Status            string `json:"status"`
	TransactionIndex  string `json:"transactionIndex"`
	BlockNum          string `json:"blockNumber"`
	GasUsed           string `json:"gasUsed"`
	Root              string `json:"string"`
	TransactionHash   string `json:"transactionHash"`
}

func getEthTx(hash string) (Tx *EthTxInfo, err error) {
	service := beego.AppConfig.String("ethurl")
	client := jsonrpc.NewClient(service)
	args := [1]string{hash}
	var resp EthTxInfo
	err1 := client.CallFor(&resp, "eth_getTransactionByHash", args)
	if err1 != nil {
		beego.Error(err1)
		return nil, err1
	}

	if resp.Hash == "" {
		beego.Error("no getEthTx result!", "hash", hash)
		return nil, nil

	}
	return &resp, nil
}

func getEthTxReceipt(hash string) (txRec *EthTxReceipt, err error) {
	service := beego.AppConfig.String("ethurl")
	client := jsonrpc.NewClient(service)
	args := [1]string{hash}
	var resp EthTxReceipt
	err1 := client.CallFor(&resp, "eth_getTransactionReceipt", args)
	if err1 != nil {
		beego.Error(err)
		return nil, err
	}
	if resp.TransactionHash == "" {
		beego.Error("no getEthReceipt result!", "hash", hash)
		return nil, nil

	}
	return &resp, nil
}

//GetEthMinerFee :get eth miner fee via tx hash
func GetEthMinerFee(hash string) (fee int64, err error) {
	tx, err := getEthTx(hash)
	if err != nil {
		return 0, err
	}
	txrec, err1 := getEthTxReceipt(hash)
	if err1 != nil {
		return 0, err
	}
	if tx == nil || txrec == nil {
		return 0, nil
	}
	gasPrice, err2 := strconv.ParseInt(tx.GasPrice[2:], 16, 64)
	if err2 != nil {
		beego.Error(err2)
		return 0, err
	}
	gasUsed, err3 := strconv.ParseInt(txrec.GasUsed[2:], 16, 64)
	if err3 != nil {
		beego.Error(err3)
		return 0, err
	}
	return gasPrice * gasUsed, nil
}

//GetEthTxFrom get  from address information of Tx
func GetEthTxFrom(hash string) (string, error) {
	tx, err := getEthTx(hash)
	if err != nil || tx == nil {
		return "", err
	}
	checkSumAddress := common.HexToAddress(tx.From)
	return checkSumAddress.Hex(), nil
}

func GetEthMinerFeeTest() {
	fee, _ := GetEthMinerFee("0xc799d542e2fe5853d1232591b4b1bca14e81fdb8e049712e046e7479163809af")
	fmt.Println("eth_fee-->", fee)
}

func ethCall(data map[string]string) (result string, err error) {
	service := beego.AppConfig.String("ethurl")
	client := jsonrpc.NewClient(service)
	args := [2]interface{}{data, "latest"}
	err1 := client.CallFor(&result, "eth_call", args)
	if err1 != nil {
		beego.Error(err1)
		return "", err1
	}

	if result == "" {
		beego.Error("no eth_call result", "data", data)
		return "", nil

	}
	return result, nil
}

//getDgatwWayAppAddress :return eth address of the given code
func getDgateWayAppAddress(appCode uint32) string {
	dgatewayAddress := beego.AppConfig.String("dgatewayContractAddress")
	funcdata := generateFuncData("getAppAddress(uint32)")
	paramdata := strconv.FormatInt(int64(appCode), 16)

	var buffer bytes.Buffer
	buffer.WriteString("0x")
	buffer.WriteString(funcdata)
	for i := 0; i < (64 - len(paramdata)); i++ {
		buffer.WriteString("0")
	}
	buffer.WriteString(paramdata)
	beego.Info("eth_call rawdata", buffer.String())
	data := map[string]string{"to": dgatewayAddress, "data": buffer.String()}
	result, err := ethCall(data)
	if err != nil {
		beego.Error(err)
		return ""
	}
	appAddress := "0x" + result[len(result)-40:]
	beego.Info("appcode:", appCode, "EthAppAddress:", appAddress)
	return appAddress
}

//getDgateWayAppName :return token symbol of given address
func getDgateWayAppSymbol(appAddress string) string {
	funcdata := generateFuncData("symbol()")
	data := map[string]string{"to": appAddress, "data": "0x" + funcdata}
	result, err := ethCall(data)
	if err != nil {
		beego.Error(err)
		panic(err)
	}
	resStr := strings.Trim(result[130:], "0")
	byteArray, _ := hex.DecodeString(resStr)
	res := string(byteArray[:])
	beego.Info("GetDgateWayAppName:", res)
	return res
}

//getDgateWayAppDecimal :return token decimals of given addredd
func getDgateWayAppDecimal(appAddress string) int {
	funcdata := generateFuncData("decimals()")
	data := map[string]string{"to": appAddress, "data": "0x" + funcdata}
	result, err := ethCall(data)
	if err != nil {
		beego.Error(err)
		panic(err)
	}
	res, _ := strconv.ParseInt(result[2:], 16, 0)
	beego.Info("GetDgateWayAppDecimal:", res)
	return int(res)
}

//GetDgateWayAppSymbolAndDecimal return symbol and address of given eth app code
func GetDgateWayAppSymbolAndDecimal(appCode uint32) (string, int) {
	appAddress := getDgateWayAppAddress(appCode)
	symbol := getDgateWayAppSymbol(appAddress)
	decimals := getDgateWayAppDecimal(appAddress)
	return symbol, decimals
}

func generateFuncData(funcName string) string {
	input := []byte(funcName)
	rawSha3 := sha3.NewKeccak256()
	rawSha3.Reset()
	rawSha3.Write(input)
	rawSha3Output := rawSha3.Sum(nil)
	result := hex.EncodeToString(rawSha3Output)
	return result[:8]
}
