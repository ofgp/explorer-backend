package chainapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//读不到beego配置,测试无效
func TestGetTokenCurrency(t *testing.T) {
	_, err := GetTokenCurrency("BCH", "CNY")
	assert.NoError(t, err)
}

func TestGetTokenCurrencyPrice(t *testing.T) {
	res := GetTokenCurrencyPrice("BCH", "CNY")
	assert.Equal(t, float64(0), res)
}
