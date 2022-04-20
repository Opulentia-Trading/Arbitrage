package platformListener

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type BinanceHandler struct {
	*cexHandler
}

func NewBinanceHandler() *BinanceHandler {
	binanceBaseUrl := "https://api.binance.com"
	binanceApiKey := ""
	binanceEndpoints := cexEndpointIdx{
		apiTest:        "/api/v3/ping",
		tickerPriceAll: "/api/v3/ticker/price",
	}

	cexHandlerInst := newCEXHandler(binanceBaseUrl, binanceApiKey, &binanceEndpoints)
	return &BinanceHandler{cexHandlerInst}
}

// Implements the ListenerHandler interface
func (b *BinanceHandler) TestConnection() {
	url := b.baseUrl + b.endpoints.apiTest
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// Free up the response body after the function completes
	defer resp.Body.Close()

	// Extract the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

// Implements the ListenerHandler interface
func (b *BinanceHandler) FetchTickerPriceAll() []TickerPrice {
	url := b.baseUrl + b.endpoints.tickerPriceAll

	// TODO: Switch to the http client and include a timeout
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.m1ailtazu
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	// Free up the response body after the function completes
	defer resp.Body.Close()

	// TODO: Check additional fields in response
	// 	- Status code: resp.StatusCode
	//	- Headers: resp.Header

	var result []TickerPrice

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// Binance Error Handling
// ======================
// - Client errors (4XX)
// 		+ Response header
// 		+ Error code
// 		+ Error message
// 		+ Request config
//
// - Server errors (5XX)
//
// Rate Limits
// -----------
// HTTP 403: WAF Limit (Web Application Firewall) violated
// HTTP 429: Breaking request rate limit, cooldown to comply with rate limits
// HTTP 418: IP has been autobanned after ignoring prior 429 codes
// The "Retry-After" header gives info of cooldown period after a HTTP 429 or 418 response
//
// /api/* endpoints have 1200 weight per min (20 sec)
// "X-MBX-USED-WEIGHT-(intervalNum)(intervalLetter)" header shows weight usage
//
// Websocket does not count towards request rate limit
