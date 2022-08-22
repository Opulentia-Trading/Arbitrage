package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Opulentia-Trading/Arbitrage/env"
	"github.com/Opulentia-Trading/Arbitrage/platformListener"
	"github.com/Opulentia-Trading/Arbitrage/platformListener/binanceHandler"
	"github.com/Opulentia-Trading/Arbitrage/platformListener/uniswapV2Handler"
	"github.com/Opulentia-Trading/Arbitrage/platformListener/uniswapV3Handler"
	"github.com/Opulentia-Trading/Arbitrage/util"
)

func main() {
	// Load the env variables
	dirname, err := util.CurDirname()
	if err != nil {
		log.Fatal(err)
	}

	dotenvPath := filepath.Join(dirname, "../../env/.env")
	env.Load_env(dotenvPath)

	// Test the Binance listener
	binanceHandler := binanceHandler.NewBinanceHandler()
	binanceListener := platformListener.NewListener("Binance", binanceHandler)

	fmt.Printf("========== %v Listener ==========\n", binanceListener.PlatformName)
	binanceListener.TestConnection()
	tickerPricesBinance := binanceListener.FetchTickerPriceAll()
	fmt.Println(tickerPricesBinance[0:5])
	tickerPriceBinance := binanceListener.FetchTickerPrice("LINK", "ETH")
	fmt.Println(tickerPriceBinance)

	// Test the UniswapV3 listener
	uniswapV3Handler := uniswapV3Handler.NewUniswapV3Handler()
	uniswapV3Listener := platformListener.NewListener("UniswapV3", uniswapV3Handler)

	fmt.Printf("\n========== %v Listener ==========\n", uniswapV3Listener.PlatformName)
	uniswapV3Listener.TestConnection()
	tickerPricesUniV3 := uniswapV3Listener.FetchTickerPriceAll()
	fmt.Println(tickerPricesUniV3)
	tickerPriceUniV3 := uniswapV3Listener.FetchTickerPrice("LINK", "ETH")
	fmt.Println(tickerPriceUniV3)

	// Test the UniswapV2 listener
	uniswapV2Handler := uniswapV2Handler.NewUniswapV2Handler()
	uniswapV2Listener := platformListener.NewListener("UniswapV2", uniswapV2Handler)

	fmt.Printf("\n========== %v Listener ==========\n", uniswapV2Listener.PlatformName)
	uniswapV2Listener.TestConnection()
	tickerPricesUniV2 := uniswapV2Listener.FetchTickerPriceAll()
	fmt.Println(tickerPricesUniV2)
	tickerPriceUniV2 := uniswapV2Listener.FetchTickerPrice("LINK", "ETH")
	fmt.Println(tickerPriceUniV2)
	fmt.Println()
	pairInfoUniV2 := uniswapV2Handler.FetchPairInfo("LINK", "ETH")
	log.Println(util.PrettyPrint(pairInfoUniV2))

	// Test Eth Gas Estimation
	fmt.Printf("\n========== Ethereum Gas Estimation ==========\n")
	platformListener.RunEthGasTests()
}
