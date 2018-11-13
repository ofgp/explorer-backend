package chainapi

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func init() {
	beego.LoadAppConfig("ini", "../conf/app.conf")
	fmt.Println("111")
}
func TestGetTokenCurrency(t *testing.T) {
	_, err := GetTokenCurrency("BCH")
	assert.NoError(t, err)
}

func TestGetTokenCurrencyPrice(t *testing.T) {
	_, err := GetTokenCurrencyPrice("BCH", "CNY")
	assert.NoError(t, err)
}

func TestGetExChangeRate(t *testing.T) {
	res, err := GetExChangeRate()
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestGetKey(t *testing.T) {
	key, _, _ := GetKey()
	assert.NotEqual(t, "", key)
}

func TestUSDToCNY(t *testing.T) {
	res, _ := USDToCNY()
	assert.NotEqual(t, 0, res)
}

func TestCNYToUSD(t *testing.T) {
	res, _ := CNYToUSD()
	assert.NotEqual(t, 0, res)
}
