package controllers

import (
	"dgatewayWebBrowser/models"

	"github.com/astaxie/beego"
)

type NodeListController struct {
	beego.Controller
}

func (nl *NodeListController) Get() {
	nodeList, err := models.GetNodeList()
	if err != nil {
		nl.Data["json"] = map[string]interface{}{"code": -1, "msg": "inner error"}
		nl.ServeJSON()
	}
	nl.Data["json"] = nodeList
	nl.ServeJSON()

}
