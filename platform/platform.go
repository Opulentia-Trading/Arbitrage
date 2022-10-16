package platform

import (
	"fmt"
	"strings"

	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/Opulentia-Trading/Arbitrage/platform/cexHandler/binanceHandler"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler/uniswapV2Handler"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler/uniswapV3Handler"
)

// Interface defining required methods for a DEX/CEX platform
type Platform interface {
	GetExchangeInfo() *models.Exchange
	TestConnection() (string, error)
	FetchTickerInfoAll() ([]models.TickerInfo, error)
	FetchTickerInfo(base string, quote string) (models.TickerInfo, error)
	ExecuteOrder(order models.Order) error
	String() string
}

func GetPlatform(platformName string) (Platform, error) {
	platformName = strings.ToLower(platformName)

	switch platformName {
	case binanceHandler.PlatformName:
		return binanceHandler.NewBinanceHandler(), nil
	case uniswapV2Handler.PlatformName:
		handler, err := uniswapV2Handler.NewUniswapV2Handler()
		if err != nil {
			return nil, err
		}
		return handler, nil
	case uniswapV3Handler.PlatformName:
		handler, err := uniswapV3Handler.NewUniswapV3Handler()
		if err != nil {
			return nil, err
		}
		return handler, nil
	default:
		return nil, fmt.Errorf("unknown platform: %v", platformName)
	}
}
