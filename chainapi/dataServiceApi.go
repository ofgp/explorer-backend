package chainapi

import (
	"dgatewayWebBrowser/datastruct"
	"fmt"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/httplib"
)

//获取代币对应法币价格
func GetTokenCurrency(target, unit string) (*datastruct.TokenCurrencyResp, error) {
	req := httplib.Post(beego.AppConfig.String("dataurl") + "/api/custom/current/price_to_currency")
	req, err := req.JSONBody(map[string]string{"target": target, "unit": unit})
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	resp := datastruct.TokenCurrencyResp{}
	err1 := req.ToJSON(&resp)
	if err != nil {
		beego.Error(err1)
		return nil, err1
	}
	fmt.Println(&resp)
	return &resp, nil
}

func GetTokenCurrencyPrice(target, unit string) float64 {
	resp, err := GetTokenCurrency(target, unit)
	if err != nil {
		return 0
	}
	return resp.Price
}
