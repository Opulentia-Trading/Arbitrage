package main

import (
	"fmt"

	"example.com/arbitrage/platformListener"
)

func main() {
	// Setup the listener
	var binanceTickers []string
	binanceListener := platformListener.NewListener("binance", binanceTickers)

	// Setup a handler for the listener
	binanceBaseUrl := "https://api.binance.com"
	binanceApiKey := ""
	binanceEndpoints := platformListener.CEXEndpointIdx{
		ApiTest:        "/api/v3/ping",
		TickerPriceAll: "/api/v3/ticker/price",
	}

	binanceHandler := platformListener.NewCEXHandler(binanceBaseUrl, binanceApiKey, &binanceEndpoints)

	// Bind the handler to the listener
	binanceListener.BindHandler(binanceHandler)

	binanceListener.PingTest()
	tickerPrices := binanceListener.FetchTickerPriceAll()
	fmt.Println(tickerPrices)
}
