// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dgatewayWebBrowser/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

/*
insert into dgateway_block (height, hash, pre_id, tx_cnt, time, size, created_used, miner)
	values (101, "0xhash99", "0xhash98", 20, 123456790, 50, 10, "miner1");
*/
func init() {
	allowcors, err := beego.AppConfig.Bool("allowcors")
	if err != nil {
		panic(err)

	}
	if allowcors {
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			//AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			AllowCredentials: true,
			AllowOrigins:     []string{beego.AppConfig.String("alloworigin")},
		}))
	}
	//首页info, tx, block的websocket连接
	beego.Router("/info/ws", &controllers.InfoWsController{})
	//首页表单
	beego.Router("/info/txdata", &controllers.InfoTableController{})
	//首页搜索
	beego.Router("/info/search", &controllers.InfoSearchController{})
	//首页总览
	beego.Router("/info/overview", &controllers.InfoOverviewController{})
	//获取最新区块高度
	beego.Router("/block/currentblock", &controllers.BlockCurrentController{})
	//获取所有节点信息
	beego.Router("/node/list", &controllers.NodeListController{})
	//获取区块列表
	beego.Router("/block/list", &controllers.BlockListController{})
	//获取/搜索区块详情
	beego.Router("/block/detail", &controllers.BlockDetailController{})
	//获取交易列表
	beego.Router("/tranx/list", &controllers.TranxListController{})
	//获取交易详情
	beego.Router("/tranx/detail", &controllers.TranxDetailController{})
	//批量发送比特币
	beego.Router("/tool/batchSendbitcoin", &controllers.BitcoinToolsController{})
}
