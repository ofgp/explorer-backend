package chainapi

import (
	"dgatewayWebBrowser/datastruct"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/httplib"
)

//GetTokenCurrency query from cgiserver
func GetTokenCurrency(target string) (*datastruct.TokenCurrencyResp, error) {
	req := httplib.Post(beego.AppConfig.String("cgiurl") + "/market/showTokenPrice")
	key, tokenCode, err := GetKey()
	fmt.Println(key, tokenCode)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	req = req.Header("TOKEN-CODE", tokenCode)
	req = req.Header("x-auth-token", key)

	reqParams := map[string]interface{}{"tokenName": target}
	finalParams, err := DoAesEncrypt(reqParams, key)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	req, err = req.JSONBody(finalParams)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	commonResp := datastruct.DataCommonResp{}
	err = req.ToJSON(&commonResp)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	finalResp := datastruct.TokenCurrencyResp{}
	decryptRes, err := DoAesDecrypt(commonResp.SData, key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(decryptRes, &finalResp)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if finalResp.Code != 1 {
		err = errors.New(finalResp.Msg)
		beego.Error(err)
		return nil, err
	}
	return &finalResp, nil
}

//GetExChangeRate return exchange rate of different currency
func GetExChangeRate() (*datastruct.ExchangeRateResp, error) {
	req := httplib.Post(beego.AppConfig.String("cgiurl") + "/mobile/transaction/getAllRate")
	resp := datastruct.ExchangeRateResp{}
	err := req.ToJSON(&resp)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	if resp.Code != 1 {
		err = errors.New(resp.Msg)
		beego.Error(err)
		return nil, err
	}
	return &resp, nil
}

//GetKey get request crypto key of cgiserve request
func GetKey() (string, string, error) {
	req := httplib.Post(beego.AppConfig.String("cgiurl") + "/mobile/wallet/getKey")
	resp, err := req.DoRequest()
	if err != nil {
		beego.Error(err)
		return "", "", err
	}
	tokenCode := resp.Header.Get("TOKEN-CODE")
	out := datastruct.DataCryptoKeyResponse{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		beego.Error(err)
		return "", "", err
	}
	if out.Code != 1 {
		err = errors.New(out.Msg)
		beego.Error(err)
		return "", "", err
	}
	return out.Data, tokenCode, nil
}

func GetTokenCurrencyPrice(target, unit string) (float64, error) {
	resp, err := GetTokenCurrency(target)
	if err != nil {
		return 0, err
	}
	beego.Info("token: ", target, "\nprice: ", resp.Data.CurrentTokenPrice)
	if unit == "CNY" {
		rate, err := USDToCNY()
		if err != nil {
			return 0, err
		}
		return resp.Data.CurrentTokenPrice * rate, nil
	} else if unit == "USD" {
		return resp.Data.CurrentTokenPrice, nil
	}
	return 0, errors.New("Unsupport currency type !")
}

//CNYToUSD CNYToUSD exchange
func USDToCNY() (float64, error) {
	resp, err := GetExChangeRate()
	if err != nil {
		return 0, err
	}
	out := float64(0)
	for _, d := range resp.Data {
		if d.Name == "CNY" {
			out = d.Rate
		}
	}
	return out, nil
}

func CNYToUSD() (float64, error) {
	resp, err := GetExChangeRate()
	if err != nil {
		return 0, err
	}
	out := float64(0)
	for _, d := range resp.Data {
		if d.Name == "CNY" {
			out = 1 / d.Rate
		}
	}
	return out, nil
}
