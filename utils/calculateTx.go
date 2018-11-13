package utils

import (
	"fmt"
	"time"

	"dgatewayWebBrowser/chainapi"
	"dgatewayWebBrowser/dboperation"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

//CalculateTxDataDaily :Get tx statistics of last hour
func CalculateTxDataDaily() {
	// calc every 10 min
	tk := toolbox.NewTask("txdailydata", "30 */10 * * * *", calcTxstatics)
	toolbox.AddTask("tk", tk)
	toolbox.StartTask()
	defer toolbox.StopTask()
}

//amount的单位锚定,根据token_code来确定
var calcTxstatics = toolbox.TaskFunc(func() error {
	//获取昨日开始和结束时间戳
	beego.Info("calculate tx statistics every 10 minutes start!!!")
	/*
		timeStr := time.Now().Format("2006-01-02")
		timeBase, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
		//包含时区信息
		startTime := timeBase.AddDate(0, 0, -1).Unix()
		endTime := timeBase.Unix()
	*/

	//change to calc every 10 min.
	timeNow := time.Now()
	startTime, endTime := getTenMinuteDuration(timeNow)
	beego.Info("calc starttime: ", startTime, "calc endtime: ", endTime)

	distinctToken, err := dboperation.GetDailyDistinctToken(startTime, endTime)
	if err != nil {
		return err
	}

	for _, token := range distinctToken {
		fmt.Println("token:", token)
		if token == "" {
			continue
		}
		tokenInfo, err := dboperation.GetTokenInfoBySymbol(token)
		fmt.Printf("token_info, %+v", tokenInfo)
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
		currencyAmount := float64(0)
		if tokenInfo.Chain == "xin" && tokenInfo.Symbol == "XIN" {

			USDToCNYPrice, err := chainapi.USDToCNY()
			if err != nil {
				beego.Critical("get exchange rate failed", err)
			}
			//eos XIN coin's exchange rate: 1USD: 1000XIN
			currencyAmount = float64(tokenAmount) / 1000 * USDToCNYPrice
		} else {
			price, err := chainapi.GetTokenCurrencyPrice(currencySymbol, "CNY")
			if err != nil {
				beego.Critical("get token price failed", err)
			}
			currencyAmount = price * float64(tokenAmount) / float64(power(currencyDecimals))
		}
		//新增数据
		time := time.Unix(startTime, 0)
		fmt.Println("token calc:", time, tokenAmount, count, tokenSymbol, currencyAmount)
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

func getTenMinuteDuration(timeNow time.Time) (start, end int64) {
	minute := (timeNow.Minute() / 10 * 10)
	endTime := time.Date(
		timeNow.Year(),
		timeNow.Month(),
		timeNow.Day(),
		timeNow.Hour(),
		minute,
		0,
		0,
		time.Local,
	)
	startTime := endTime.Add(-10 * time.Minute)
	return startTime.Unix(), endTime.Unix()
}
