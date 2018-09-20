package controllers

import (
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/models"
	"encoding/json"

	"github.com/astaxie/beego"
)

type BitcoinToolsController struct {
	beego.Controller
}

func (bc *BitcoinToolsController) Post() {
	var req models.BatchSendRequest
	err := json.Unmarshal(bc.Ctx.Input.RequestBody, &req)
	if err != nil {
		bc.Data["json"] = datastruct.BatchSendResp{Code: -1, Msg: "param error"}
		bc.ServeJSON()
	}

	res := models.BatchSendBitcoin(&req)
	bc.Data["json"] = res
	bc.ServeJSON()
}
