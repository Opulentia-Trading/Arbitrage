package cexHandler

import "github.com/Opulentia-Trading/Arbitrage/models"

type CexEndpointIdx struct {
	ApiTest        string
	TickerPriceAll string
	TickerPrice    string
}

type CexHandler struct {
	ExchangeInfo *models.Exchange
	BaseUrl      string
	ApiKey       string
	Endpoints    *CexEndpointIdx
}

func NewCEXHandler(exchangeInfo *models.Exchange, baseUrl string, apiKey string, endpoints *CexEndpointIdx) *CexHandler {
	return &CexHandler{
		ExchangeInfo: exchangeInfo,
		BaseUrl:      baseUrl,
		ApiKey:       apiKey,
		Endpoints:    endpoints,
	}
}
