package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Opulentia-Trading/Arbitrage/env"
	"github.com/Opulentia-Trading/Arbitrage/platformListener"
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
	binanceHandler := platformListener.NewBinanceHandler()
	binanceListener := platformListener.NewListener("Binance", binanceHandler)

	fmt.Printf("========== %v Listener ==========\n", binanceListener.PlatformName)
	binanceListener.TestConnection()
	tickerPricesBinance := binanceListener.FetchTickerPriceAll()
	fmt.Println(tickerPricesBinance[0:5])
	tickerPriceBinance := binanceListener.FetchTickerPrice("LINK", "ETH")
	fmt.Println(tickerPriceBinance)

	// Test the UniswapV3 listener
	uniswapV3Handler := platformListener.NewUniswapV3Handler()
	uniswapV3Listener := platformListener.NewListener("UniswapV3", uniswapV3Handler)

	fmt.Printf("\n========== %v Listener ==========\n", uniswapV3Listener.PlatformName)
	uniswapV3Listener.TestConnection()
	tickerPricesUniswap := uniswapV3Listener.FetchTickerPriceAll()
	fmt.Println(tickerPricesUniswap)
	tickerPriceUniswap := uniswapV3Listener.FetchTickerPrice("LINK", "ETH")
	fmt.Println(tickerPriceUniswap)
}
