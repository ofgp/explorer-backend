package utils

import "dgatewayWebBrowser/chainapi"

// CNYToUSD CNYToUSD exchange, but not a good practice.
func CNYToUSD() float64 {
	USDPrice := chainapi.GetTokenCurrencyPrice("BTC", "USD")
	CNYPrice := chainapi.GetTokenCurrencyPrice("BTC", "CNY")
	return USDPrice / CNYPrice
}

func USDToCNY() float64 {
	USDPrice := chainapi.GetTokenCurrencyPrice("BTC", "USD")
	CNYPrice := chainapi.GetTokenCurrencyPrice("BTC", "CNY")
	return CNYPrice / USDPrice
}
