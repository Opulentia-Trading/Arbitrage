package platformListener

type CexEndpointIdx struct {
	ApiTest        string
	TickerPriceAll string
	TickerPrice    string
}

type CexHandler struct {
	BaseUrl   string
	ApiKey    string
	Endpoints *CexEndpointIdx
}

func NewCEXHandler(baseUrl string, apiKey string, endpoints *CexEndpointIdx) *CexHandler {
	return &CexHandler{baseUrl, apiKey, endpoints}
}
