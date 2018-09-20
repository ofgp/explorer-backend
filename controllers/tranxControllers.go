package controllers

import (
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/models"

	"github.com/astaxie/beego"
)

type TranxListController struct {
	beego.Controller
}

func (tl *TranxListController) Get() {
	page, err := tl.GetInt("page")
	if err != nil {
		beego.Error(err)
		tl.Ctx.WriteString("no page params")
		return
	}
	pageSize, err1 := tl.GetInt("page_size")
	if err1 != nil {
		beego.Error(err)
		tl.Ctx.WriteString("no page_size params")
		return
	}
	if page <= 0 || pageSize <= 0 {
		tl.Ctx.WriteString("illegal page or page_size")
		return
	}
	search := tl.GetString("search")
	resp, err := models.GetTxList(page, pageSize, search)
	if err != nil {
		tl.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		tl.ServeJSON()
	}
	tl.Data["json"] = resp
	tl.ServeJSON()
}

type TranxDetailController struct {
	beego.Controller
}

func (td *TranxDetailController) Get() {
	dgwHash := td.GetString("dgw_hash")
	resp, err := models.GetTxByHash(dgwHash)
	if err != nil {
		td.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		td.ServeJSON()
	}
	td.Data["json"] = resp
	td.ServeJSON()
}
