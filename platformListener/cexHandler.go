package platformListener

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type CEXEndpointIdx struct {
	ApiTest        string
	TickerPriceAll string
}

type CEXHandler struct {
	baseUrl   string
	apiKey    string
	endpoints *CEXEndpointIdx
}

func NewCEXHandler(baseUrl string, apiKey string, endpoints *CEXEndpointIdx) *CEXHandler {
	return &CEXHandler{baseUrl, apiKey, endpoints}
}

func (c *CEXHandler) PingTest() {
	url := c.baseUrl + c.endpoints.ApiTest
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

func (c *CEXHandler) FetchTickerPriceAll() []TickerPrice {
	url := c.baseUrl + c.endpoints.TickerPriceAll

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
