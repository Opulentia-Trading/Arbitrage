package validator

import (
	"github.com/Opulentia-Trading/Arbitrage/platformListener"
)

type TickerPrice = platformListener.TickerPrice

type PlatformInfo struct {
	TYPE string
	NAME string
}
type ArbResult struct {
	SYMBOL   		string
	PLATFORM 		*PlatformInfo
	ACTION        	string
	AMOUNT        	int
	TRANSACTIONID 	string
	NEXT          	*ArbResult
}

func main(platform1 TickerPrice, platform2 TickerPrice) ArbResult{
	var result ArbResult
	// Get Flashbot contract

	// Execute contract getProfit function between both pair pools



	return result
}