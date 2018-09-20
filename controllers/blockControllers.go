package controllers

import (
	"dgatewayWebBrowser/datastruct"
	"dgatewayWebBrowser/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BlockListController struct {
	beego.Controller
}

//Get :获取指定范围区块信息
// @Param	page_size		path 	string	true
// @Param	page		path 	string	true
func (bl *BlockListController) Get() {
	page, err := bl.GetInt("page")
	if err != nil {
		beego.Error(err)
		bl.Ctx.WriteString("no page params")
		return
	}
	pageSize, err1 := bl.GetInt("page_size")
	if err1 != nil {
		beego.Error(err)
		bl.Ctx.WriteString("no page_size params")
		return
	}
	if page <= 0 || pageSize <= 0 {
		bl.Ctx.WriteString("illegal page or page_size")
		return
	}
	search := bl.GetString("search")
	res, err := models.GetBlockList(page, pageSize, search)
	if err != nil {
		bl.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		bl.ServeJSON()
	}
	bl.Data["json"] = res
	bl.ServeJSON()

}

type BlockDetailController struct {
	beego.Controller
}

func (bd *BlockDetailController) Get() {
	height := bd.GetString("height")
	if height == "" {
		bd.Ctx.WriteString("Invalid height params")
		return
	}
	res, err := models.GetBlockDetail(height)
	//no row found
	if err == orm.ErrNoRows {
		bd.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "No row found"}
		bd.ServeJSON()
	}
	if err != nil {
		bd.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		bd.ServeJSON()
	}
	bd.Data["json"] = res
	bd.ServeJSON()
}

type BlockCurrentController struct {
	beego.Controller
}

func (bc *BlockCurrentController) Get() {
	res, err := models.GetCurrentBlock()
	if err != nil {
		bc.Data["json"] = datastruct.ErrResp{Code: -1, Msg: "inner error"}
		bc.ServeJSON()
	}
	bc.Data["json"] = res
	bc.ServeJSON()
}
