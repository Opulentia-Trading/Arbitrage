package main

import (
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/env"
	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/Opulentia-Trading/Arbitrage/platform"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler/uniswapV2Handler"
	"github.com/Opulentia-Trading/Arbitrage/util"
)

const TickerLimit = 5

func platformTest(platformName string, base string, quote string) {
	platform, err := platform.GetPlatform(platformName)
	if err != nil {
		panic(err)
	}

	platformInfo := platform.GetExchangeInfo()
	fmt.Printf("========== %v platform test ==========\n", platformInfo.Name)
	fmt.Println("Name:", platformInfo.Name)
	fmt.Println("Type:", platformInfo.Type)

	connTest, err := platform.TestConnection()
	if err != nil {
		panic(err)
	}
	fmt.Println("\n+--------- Connection Test ---------+")
	fmt.Println(connTest)

	tickerInfoAll, err := platform.FetchTickerInfoAll()
	if err != nil {
		panic(err)
	}
	if len(tickerInfoAll) > TickerLimit {
		tickerInfoAll = tickerInfoAll[:TickerLimit]
	}
	fmt.Println("\n+--------- Fetch Tickers ---------+")
	fmt.Printf("limit=%v\n", TickerLimit)
	fmt.Println(util.PrettyPrint(tickerInfoAll))

	tickerInfo, err := platform.FetchTickerInfo(base, quote)
	if err != nil {
		panic(err)
	}
	fmt.Println("\n+--------- Fetch Ticker ---------+")
	fmt.Println(util.PrettyPrint(tickerInfo))

	if u, ok := platform.(*uniswapV2Handler.UniswapV2Handler); ok {
		reserves, err := u.FetchPairReserves(base, quote)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n+--------- Fetch Pair Reserves ---------+")
		fmt.Println(util.PrettyPrint(reserves))
	}

	testOrder := models.Order{
		Exchange: platformInfo,
		Base:     base,
		Quote:    quote,
		Action:   models.BuyLongSpot,
		Price:    new(big.Rat).SetFloat64(24.86),
		Quantity: big.NewInt(200),
		Deadline: time.Minute,
	}
	fmt.Println("\n+--------- Execute Order ---------+")
	platform.ExecuteOrder(testOrder)
}

func main() {
	// Load env variables
	dirname, err := util.CurDirname()
	if err != nil {
		panic(err)
	}

	dotenvPath := filepath.Join(dirname, "../../env/.env")
	env.Load_env(dotenvPath)

	platformNames := []string{"binance", "uniswap_v2", "uniswap_v3"}
	for _, platformName := range platformNames {
		platformTest(platformName, "LINK", "ETH")
		fmt.Print("\n\n\n")
	}

	fmt.Println("========== eth gas estimation ==========")
	ethHandler.RunEthGasTests()
}
