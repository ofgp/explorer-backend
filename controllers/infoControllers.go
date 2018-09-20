package controllers

import (
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type InfoWsController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

func (ic *InfoWsController) Get() {
	//ws, err := upgrader.Upgrade(ic.Ctx.ResponseWriter, ic.Ctx.Request, nil)
	isWs := false
	for conn := range ic.Ctx.Request.Header {
		if conn == "Upgrade" {
			isWs = true
			break
		}
	}
	if !isWs {
		beego.Info("wrong method")
		ic.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "wrong method"}
		ic.ServeJSON()
		return
	}
	//TODO: 切换成Upgrader
	ws, err := websocket.Upgrade(ic.Ctx.ResponseWriter, ic.Ctx.Request, nil, 1024, 1024)
	if err != nil {
		beego.Error(err)
		ic.Data["json"] = map[string]interface{}{"code": -1, "msg": "inner error", "data": map[string]string{}}
		ic.ServeJSON()
	}
	models.WsClients[ws] = true

}

//首页图表
type InfoTableController struct {
	beego.Controller
}

func (it *InfoTableController) Get() {
	resp, err := models.GetInfoTableData()
	if err != nil {
		it.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		it.ServeJSON()
	}
	it.Data["json"] = resp
	it.ServeJSON()
}

//首页搜索
type InfoSearchController struct {
	beego.Controller
}

func (is *InfoSearchController) Get() {
	searchBlock, err := is.GetInt64("search")
	//搜索区块
	if err == nil {
		search := strconv.FormatInt(searchBlock, 10)
		res, err1 := models.GetBlockList(1, 10, search)
		if err1 != nil {
			is.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
			is.ServeJSON()
			return
		}
		is.Data["json"] = res
		is.ServeJSON()
		return
	}
	searchTx := is.GetString("search")
	res, err2 := models.GetTxList(1, 10, searchTx)
	if err2 != nil {
		is.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		is.ServeJSON()
		return
	}
	is.Data["json"] = res
	is.ServeJSON()
}

//首页总览(当前高度， 节点数量， 交易数量， 区块总数）
type InfoOverviewController struct {
	beego.Controller
}

func (io *InfoOverviewController) Get() {
	data, err := models.GetInfoOverview()
	if err != nil {
		io.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner errror"}
		io.ServeJSON()
		return
	}
	io.Data["json"] = data
	io.ServeJSON()
}
