package utils

import (
	"fmt"
	"time"

	"dgatewayWebBrowser/chainapi"
	"dgatewayWebBrowser/dboperation"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

//CalculateTxDataDaily :Get tx statistics of yesterday
func CalculateTxDataDaily() {
	tk := toolbox.NewTask("txdailydata", "30 0 0 * *", calcTxstatics)
	//tk := toolbox.NewTask("txdailydata", "*/5 * * * *", calcTxstatics)
	toolbox.AddTask("tk", tk)
	toolbox.StartTask()
	defer toolbox.StopTask()

}

//amount的单位锚定,根据token_code来确定
var calcTxstatics = toolbox.TaskFunc(func() error {
	//获取昨日开始和结束时间戳
	beego.Info("calculate daily statistics start!!!")
	timeStr := time.Now().Format("2006-01-02")
	timeBase, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	//包含时区信息
	startTime := timeBase.AddDate(0, 0, -1).Unix()
	endTime := timeBase.Unix()

	distinctToken, err := dboperation.GetDailyDistinctToken(startTime, endTime)
	if err != nil {
		return err
	}

	for _, token := range distinctToken {
		if token == "" {
			continue
		}
		tokenInfo, err := dboperation.GetTokenInfoBySymbol(token)
		if err != nil {
			continue
		}
		tokenSymbol := tokenInfo.Symbol

		//计算单种token下的交易总金额和交易总数
		count, tokenAmount, err := dboperation.GetDailyTokenAmountAndCountFromMysql(startTime, endTime, tokenSymbol)
		if err != nil {
			continue
		}
		//计算token对应的CNY法币金额, eth上的token需等值到相关链的token
		currencySymbol := ""
		currencyDecimals := 1
		if tokenInfo.Chain == "eth" {
			relateToken, err := dboperation.GetTokenInfo(tokenInfo.RelateChain, tokenInfo.RelateTokenCode)
			if err != nil {
				continue
			}
			currencySymbol = relateToken.Symbol
			currencyDecimals = relateToken.Decimals
		} else {

			currencySymbol = tokenInfo.Symbol
			currencyDecimals = tokenInfo.Decimals
		}
		price := chainapi.GetTokenCurrencyPrice(currencySymbol, "CNY")
		currencyAmount := price * float64(tokenAmount) / float64(power(currencyDecimals))
		//新增数据
		time := timeBase.AddDate(0, 0, -1)
		fmt.Println(currencySymbol, count, price, tokenAmount, currencyAmount)
		err = dboperation.StoreTxStatisticToMysql(time, tokenAmount, count, tokenSymbol, currencyAmount)
		if err != nil {
			return err
		}
	}
	return nil
})

func power(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
