package main

import (
	"dgatewayWebBrowser/chainwatcher"
	_ "dgatewayWebBrowser/routers"
	"dgatewayWebBrowser/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.LoadAppConfig("ini", "conf/app.conf")

	// 跨域设置
	allowOrigins := beego.AppConfig.Strings("allow_origins")
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		AllowOrigins:     allowOrigins,
	}))

	go chainwatcher.StartWatch() //抓取链上区块信息
	utils.CalculateTxDataDaily() //每日计算交易数据
	//设置logs输出
	beego.SetLogger("file", `{"filename":"logs/app.log"}`)
	beego.Run()
}
