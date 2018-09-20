package main

import (
	"dgatewayWebBrowser/chainwatcher"
	_ "dgatewayWebBrowser/routers"
	"dgatewayWebBrowser/utils"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.LoadAppConfig("ini", "conf/app.conf")
	go chainwatcher.StartWatch() //抓取链上区块信息
	utils.CalculateTxDataDaily() //每日计算交易数据
	//设置logs输出
	beego.SetLogger("file", `{"filename":"logs/app.log"}`)
	beego.Run()
}
