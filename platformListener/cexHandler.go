package platformListener

type cexEndpointIdx struct {
	apiTest        string
	tickerPriceAll string
	tickerPrice    string
}

type cexHandler struct {
	baseUrl   string
	apiKey    string
	endpoints *cexEndpointIdx
}

func newCEXHandler(baseUrl string, apiKey string, endpoints *cexEndpointIdx) *cexHandler {
	return &cexHandler{baseUrl, apiKey, endpoints}
}
