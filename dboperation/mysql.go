package dboperation

import (
	"dgatewayWebBrowser/datastruct"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	idleconns, _ := beego.AppConfig.Int("mysqlmaxidleconns")
	openconns, _ := beego.AppConfig.Int("mysqlmaxopenconns")
	orm.RegisterDataBase(
		"default",
		"mysql",
		beego.AppConfig.String("mysqlurl"),
		idleconns,
		openconns,
	)
	orm.RegisterModel(
		new(datastruct.DgatewayTx),
		new(datastruct.DgatewayScanHeight),
		new(datastruct.DgatewayBlock),
		new(datastruct.DgatewayTxStatistics),
		new(datastruct.DgatewayTokenInfo),
	)
	orm.RunSyncdb("default", false, true)
}

//GetTokenInfo return tokeninfo of given condition
func GetTokenInfo(chain string, tokenCode uint32) (resp *datastruct.DgatewayTokenInfo, err error) {
	o := orm.NewOrm()
	sql := fmt.Sprintf("select * from dgateway_token_info where chain = '%s' and token_code = %d", chain, tokenCode)
	err1 := o.Raw(sql).QueryRow(&resp)
	if err1 != nil {
		beego.Error(err1, "please configure token info", "chain:", chain, "token_code:", tokenCode)
		return nil, err1
	}
	return resp, nil
}

//GetTokenInfoBySymbol return tokeninfo of given symbol
func GetTokenInfoBySymbol(token string) (resp *datastruct.DgatewayTokenInfo, err error) {
	o := orm.NewOrm()
	sql := fmt.Sprintf("select * from dgateway_token_info where symbol = '%s'", token)
	err = o.Raw(sql).QueryRow(&resp)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return resp, nil
}

//AddTokenInfo add one token info to mysql of given params
func AddTokenInfo(chain, symbol, relateChain string, decimals int, tokenCode, relateTokenCode uint32) bool {
	o := orm.NewOrm()
	tokenInfo := &datastruct.DgatewayTokenInfo{
		Chain:           chain,
		Symbol:          symbol,
		Decimals:        decimals,
		RelateChain:     relateChain,
		TokenCode:       tokenCode,
		RelateTokenCode: relateTokenCode,
	}
	//insert
	err := o.Begin()
	if err != nil {
		beego.Error(err)
		return false
	}
	_, err1 := o.Insert(tokenInfo)
	if err1 != nil {
		beego.Error(err1)
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

//Read or create scanHeight from mysql
func ScanHeightRead() (height int64, err error) {
	o := orm.NewOrm()
	scanHeight := datastruct.DgatewayScanHeight{Name: "scan_height"}
	if created, _, err := o.ReadOrCreate(&scanHeight, "Name"); err == nil {
		if created {
			return 1, nil
		} else {
			return scanHeight.Height, nil
		}

	} else {
		beego.Error(err)
		return 1, err
	}
}

//Update scanHeight to Mysql
func ScanHeightUpdate(height int64) bool {
	scanHeight := datastruct.DgatewayScanHeight{
		Name: "scan_height",
	}
	o := orm.NewOrm()
	err := o.Read(&scanHeight, "Name")
	if err != nil {
		beego.Error(err)
		return false
	}
	scanHeight.Height = height
	_, err1 := o.Update(&scanHeight)
	if err1 != nil {
		beego.Error(err)
		return false
	}
	return true
}

//save tx
func StoreTxToMysql(tx *datastruct.DgatewayTx) bool {
	o := orm.NewOrm()
	//check exist
	exist := false
	dgatewayTx := datastruct.DgatewayTx{DgwHash: tx.DgwHash}
	err := o.Read(&dgatewayTx, "DgwHash")
	if err == nil {
		exist = true
	}
	beego.Debug("exist in mysql:", exist, "dgwhash:", tx.DgwHash)
	if exist {
		//update
		tx.Id = dgatewayTx.Id
		err := o.Begin()
		if err != nil {
			beego.Error(err)
			return false
		}
		_, err1 := o.Update(tx)
		if err1 != nil {
			beego.Error(err1)
			o.Rollback()
			return false
		}
		o.Commit()
		return true
	}
	//insert
	err = o.Begin()
	if err != nil {
		beego.Error(err)
		return false
	}
	_, err1 := o.Insert(tx)
	if err1 != nil {
		beego.Error(err1)
		o.Rollback()
		return false
	}
	o.Commit()
	return true

}

func GetTxListFromMysql(limit0, limit1 int64, search string) (count int64, resp []datastruct.DgatewayTx, err error) {
	o := orm.NewOrm()
	if search != "" {
		_, err1 := o.Raw("SELECT * FROM dgateway_tx WHERE block_height = ? or dgw_hash = ? or from_tx_hash = ? or to_tx_hash = ? ORDER BY id DESC LIMIT ?, ?",
			search, search, search, search, limit0, limit1).QueryRows(&resp)
		if err1 != nil {
			beego.Error(err1)
			return 0, nil, err1
		}
	} else {
		_, err2 := o.Raw("SELECT * FROM dgateway_tx ORDER BY id DESC LIMIT ?, ?", limit0, limit1).QueryRows(&resp)
		if err2 != nil {
			beego.Error(err2)
			return 0, nil, err2
		}
	}
	count = GetTxCount(search)
	return count, resp, nil
}

func GetTxFromMysql(dgwHash string) (*datastruct.DgatewayTx, error) {
	o := orm.NewOrm()
	tx := datastruct.DgatewayTx{DgwHash: dgwHash}
	err := o.Read(&tx, "dgw_hash")
	if err == orm.ErrNoRows {
		beego.Debug(err)
		return &tx, nil
	}
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return &tx, nil
}

func GetTxCount(search string) int64 {
	o := orm.NewOrm()
	dgatewayTx := new(datastruct.DgatewayTx)
	if search == "" {
		cnt, err1 := o.QueryTable(dgatewayTx).Count()
		if err1 != nil {
			beego.Error(err1)
		}
		return cnt
	} else {
		cond := orm.NewCondition()
		cond1 := cond.Or("block_height", search).Or("dgw_hash", search).Or("from_tx_hash", search).Or("to_tx_hash", search)
		cnt, err2 := o.QueryTable(dgatewayTx).SetCond(cond1).Count()
		if err2 != nil {
			beego.Error(err2)
		}
		return cnt
	}
}

func StoreBlockToMysql(block *datastruct.DgatewayBlock) bool {
	o := orm.NewOrm()
	//check exist
	exist := false
	dgatewayBlock := datastruct.DgatewayBlock{Height: block.Height}
	err := o.Read(&dgatewayBlock, "Height")
	if err == nil {
		exist = true
	}
	if exist {
		//update
		block.Id = dgatewayBlock.Id
		err := o.Begin()
		if err != nil {
			beego.Error(err)
			return false
		}
		_, err1 := o.Update(block)
		if err1 != nil {
			beego.Error(err1)
			o.Rollback()
			return false
		}
		o.Commit()
		return true
	}
	//insert
	err = o.Begin()
	if err != nil {
		beego.Error(err)
		return false
	}
	_, err1 := o.Insert(block)
	if err1 != nil {
		beego.Error(err1)
		o.Rollback()
		return false
	}
	o.Commit()
	return true
}

func GetBlockCount(search string) int64 {
	o := orm.NewOrm()
	dgatewayblock := new(datastruct.DgatewayBlock)
	if search == "" {
		cnt, err := o.QueryTable(dgatewayblock).Count()
		if err != nil {
			beego.Error(err)
		}
		return cnt
	} else {
		cnt, err1 := o.QueryTable(dgatewayblock).Filter("height", search).Count()
		if err1 != nil {
			beego.Error(err1)
		}
		return cnt
	}
}

//获取最新一个区块
func GetLatestBlock() (*datastruct.DgatewayBlock, error) {
	o := orm.NewOrm()
	var block datastruct.DgatewayBlock
	err := o.Raw("SELECT * FROM dgateway_block order by id desc limit 1").QueryRow(&block)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return &block, nil
}

//根据传入的limit0和Limit1返回对应block列表
func GetBlockListFromMysql(limit0, limit1 int64, search string) (count int64, resp []datastruct.DgatewayBlock, err error) {
	o := orm.NewOrm()
	if search != "" {
		_, err1 := o.Raw("SELECT * FROM dgateway_block WHERE height = ? or hash = ?", search, search).QueryRows(&resp)
		if err1 != nil {
			beego.Error(err)
			return 0, nil, err1
		}
	} else {
		_, err2 := o.Raw("SELECT *  FROM dgateway_block ORDER BY id DESC LIMIT ?,?", limit0, limit1).QueryRows(&resp)
		if err != nil {
			beego.Error(err)
			return 0, nil, err2
		}
	}
	count = GetBlockCount(search)
	return count, resp, nil

}

//根据block height获取区块信息
func GetBlockDataFromMysql(height string) (*datastruct.DgatewayBlock, error) {
	o := orm.NewOrm()
	var block datastruct.DgatewayBlock
	err := o.Raw("SELECT * FROM dgateway_block where concat(height) = ? or hash = ?", height, height).QueryRow(&block)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return &block, nil
}

func GetBlockTxFromMysql(height string) ([]datastruct.DgatewayTx, error) {
	o := orm.NewOrm()
	var txs []datastruct.DgatewayTx
	_, err := o.Raw("SELECT * FROM dgateway_tx where block = ? or concat(block_height) = ?", height, height).QueryRows(&txs)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return txs, nil
}

//获取首页表单数据
func GetInfoDataFromMysql(startTime, endTime time.Time) ([]datastruct.DgatewayTxStatistics, error) {
	o := orm.NewOrm()
	var info []datastruct.DgatewayTxStatistics
	_, err := o.QueryTable("dgateway_tx_statistics").Filter("time__gte", startTime).Filter("time__lte", endTime).All(&info)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	return info, nil
}

//获取时间段内fromChain, tochain信息
func GetDailyDistinctChainFromMysql(queryCol string, startTime, endTime int64) ([]string, error) {
	o := orm.NewOrm()
	var resp []datastruct.DgatewayTx
	sql := fmt.Sprintf("SELECT DISTINCT(%s) FROM dgateway_tx WHERE time >= %d AND time < %d", queryCol, startTime, endTime)
	_, err := o.Raw(sql).QueryRows(&resp)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	var res []string
	for _, i := range resp {
		res = append(res, i.FromChain)
	}
	return res, nil
}

//获取特定FromChain下的总金额
func GetDailyAmountOfChainFromMysql(startTime, endTime int64, queryChain, chain string) (int64, error) {
	o := orm.NewOrm()
	var res int64
	sql := fmt.Sprintf("select sum(amount) from dgateway_tx where time >= %d and time < %d and %s = '%s'", startTime, endTime, queryChain, chain)
	err := o.Raw(sql).QueryRow(&res)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	return res, nil
}

//获取特定fromchain下的总数
func GetDailyCountOfChainFromMysql(startTime, endTime int64, queryChain, chain string) (int64, error) {
	o := orm.NewOrm()
	var res int64
	sql := fmt.Sprintf("select count(*) from dgateway_tx where time >= %d and time < %d and %s = '%s' ", startTime, endTime, queryChain, chain)
	err := o.Raw(sql).QueryRow(&res)
	if err != nil {
		beego.Error(err)
		return 0, err
	}
	return res, nil
}

//获取distinct token
func GetDailyDistinctToken(startTime, endTime int64) ([]string, error) {
	o := orm.NewOrm()
	var resp []datastruct.DgatewayTx
	sql := fmt.Sprintf("SELECT DISTINCT(token_symbol) FROM dgateway_tx WHERE time >= %d AND time < %d", startTime, endTime)
	_, err := o.Raw(sql).QueryRows(&resp)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	var res []string
	for _, i := range resp {
		res = append(res, i.TokenSymbol)
	}
	return res, nil
}

//获取每日特定token下面的总数及总金额
func GetDailyTokenAmountAndCountFromMysql(startTime, endTime int64, tokenSymbol string) (int64, int64, error) {
	o := orm.NewOrm()
	var count int64
	sql1 := fmt.Sprintf("select count(*) from dgateway_tx where time >= %d and time < %d and token_symbol = '%s'", startTime, endTime, tokenSymbol)
	err1 := o.Raw(sql1).QueryRow(&count)
	if err1 != nil {
		beego.Error(err1)
		return 0, 0, err1
	}
	var amount int64
	sql2 := fmt.Sprintf("select sum(amount) from dgateway_tx where time >= %d and time < %d and token_symbol = '%s'", startTime, endTime, tokenSymbol)
	err2 := o.Raw(sql2).QueryRow(&amount)
	if err2 != nil {
		beego.Error(err2)
		return 0, 0, err2
	}
	return count, amount, nil
}

//新增tx_statistics到mysql
func StoreTxStatisticToMysql(time time.Time, amount, count int64, symbol string, currencyAmount float64) error {
	o := orm.NewOrm()
	newTxStatistic := datastruct.DgatewayTxStatistics{
		Time:           time,
		Amount:         amount,
		Count:          count,
		Symbol:         symbol,
		CurrencyAmount: currencyAmount,
	}
	err := o.Begin()
	if err != nil {
		beego.Error(err)
		return err
	}
	_, err1 := o.Insert(&newTxStatistic)
	if err1 != nil {
		beego.Error(err)
		o.Rollback()
		return err1
	}
	o.Commit()
	return nil

}
