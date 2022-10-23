package binanceHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/Opulentia-Trading/Arbitrage/platform/cexHandler"
)

const (
	PlatformName      = "binance"
	httpClientTimeout = 10 * time.Second
)

var (
	// Default client does not specify a timeout
	httpClient = &http.Client{
		Timeout: httpClientTimeout,
	}

	quoteRegexOnce sync.Once
	quoteRegex     *regexp.Regexp
)

// Implements the Platform interface
type BinanceHandler struct {
	*cexHandler.CexHandler
}

func NewBinanceHandler() *BinanceHandler {
	exchangeInfo := models.Exchange{
		Type: models.Centralized,
		Name: PlatformName,
	}

	baseUrl := "https://api.binance.com"
	apiKey := ""
	endpoints := cexHandler.CexEndpointIdx{
		ApiTest:        "/api/v3/time",
		TickerPriceAll: "/api/v3/ticker/price",
		TickerPrice:    "/api/v3/ticker/price?symbol=",
	}

	cexHandlerInst := cexHandler.NewCEXHandler(&exchangeInfo, baseUrl, apiKey, &endpoints)
	initQuoteRegex(baseUrl)
	return &BinanceHandler{cexHandlerInst}
}

func (h *BinanceHandler) GetExchangeInfo() *models.Exchange {
	return h.ExchangeInfo
}

func (h *BinanceHandler) TestConnection() (string, error) {
	url := h.BaseUrl + h.Endpoints.ApiTest
	resp, err := httpClient.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body), nil
}

func (h *BinanceHandler) FetchTickerInfoAll() ([]models.TickerInfo, error) {
	url := h.BaseUrl + h.Endpoints.TickerPriceAll
	resp, err := httpClient.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var result []models.TickerInfo
	dec := json.NewDecoder(resp.Body)

	// Read array open bracket
	t, err := dec.Token()
	if err != nil {
		panic(err)
	}
	if t != json.Delim('[') {
		err := errors.New("response must be an array")
		panic(err)
	}

	// Read array contents
	for dec.More() {
		var ticker models.TickerInfo
		err := dec.Decode(&ticker)
		if err != nil {
			panic(err)
		}

		symbolSplit := quoteRegex.FindAllStringSubmatch(ticker.Symbol, -1)
		if len(symbolSplit) > 0 {
			ticker.Base = symbolSplit[0][1]
			ticker.Quote = symbolSplit[0][2]
		} else {
			ticker.Base = ticker.Symbol
			ticker.Quote = ticker.Symbol
		}

		// TODO: Use GET /sapi/v1/asset/tradeFee signed endpoint
		ticker.MakerComission = "0.001"
		ticker.TakerComission = "0.001"

		ticker.Timestamp = time.Now()
		result = append(result, ticker)
	}

	return result, nil
}

func (h *BinanceHandler) FetchTickerInfo(base string, quote string) (models.TickerInfo, error) {
	var result models.TickerInfo

	symbol := base + quote
	url := h.BaseUrl + h.Endpoints.TickerPrice + symbol
	resp, err := httpClient.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		panic(err)
	}

	result.Base = base
	result.Quote = quote

	// TODO: Use GET /sapi/v1/asset/tradeFee signed endpoint
	result.MakerComission = "0.001"
	result.TakerComission = "0.001"

	result.Timestamp = time.Now()
	return result, nil
}

func (h *BinanceHandler) ExecuteOrder(order models.Order) error {
	fmt.Printf("Executing %v/%v %v order\n", order.Base, order.Quote, order.Action.String())
	return nil
}

func (h *BinanceHandler) String() string {
	return h.ExchangeInfo.Name
}

type BinanceExchInfoResponse struct {
	Symbols []struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
	} `json:"symbols"`
}

func initQuoteRegex(baseUrl string) {
	quoteRegexOnce.Do(func() {
		// Guaranteed to run only once during program lifetime
		quoteRegexStr, err := getQuoteRegexStr(baseUrl)
		if err != nil {
			panic(err)
		}
		quoteRegex = regexp.MustCompile(quoteRegexStr)
	})
}

func getQuoteRegexStr(baseUrl string) (string, error) {
	url := baseUrl + "/api/v3/exchangeInfo"
	resp, err := httpClient.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var respData BinanceExchInfoResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&respData)
	if err != nil {
		panic(err)
	}

	quoteAssetMap := make(map[string]bool)
	var sb strings.Builder
	sb.WriteString(`^(\w+)(`)
	writeDelim := false

	for _, symbol := range respData.Symbols {
		if _, found := quoteAssetMap[symbol.QuoteAsset]; !found {
			quoteAssetMap[symbol.QuoteAsset] = true
			if writeDelim {
				sb.WriteString(`|`)
			}
			sb.WriteString(symbol.QuoteAsset)
			writeDelim = true
		}
	}

	if len(quoteAssetMap) == 0 {
		panic("could not retrieve quote assets")
	}

	sb.WriteString(`)$`)
	return sb.String(), nil
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
