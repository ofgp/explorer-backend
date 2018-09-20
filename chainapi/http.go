package chainapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
)

//基础Get请求
func httpGet(_url string, _params map[string]string) ([]byte, error) {
	baseUrl, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for k, v := range _params {
		u := &url.URL{Path: k}
		k = u.String()

		u = &url.URL{Path: v}
		v = u.String()

		params.Add(k, v)
	}

	baseUrl.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", baseUrl.String(), strings.NewReader(""))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		beego.Error(err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	httpError := HTTPError{}
	if httpError.Error.Code != 0 {
		errStr, _ := HTTPErrorTOJSON(httpError)
		beego.Error(errStr)
		return nil, errors.New(errStr)
	}
	return body, nil
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   struct {
		Code    int    `json:"code"`
		Name    string `json:"name"`
		What    string `json:"what"`
		Details []struct {
			Message    string `json:"message"`
			File       string `json:"file"`
			LineNumber int    `json:"line_number"`
			Method     string `json:"method"`
		} `json:"details"`
	} `json:"error"`
}

func HTTPErrorTOJSON(httperr HTTPError) (string, error) {

	json, err := json.Marshal(httperr)

	if err != nil {
		return "", err
	}
	return string(json), nil
}
