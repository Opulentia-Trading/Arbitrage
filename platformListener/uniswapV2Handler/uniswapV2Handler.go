package uniswapV2Handler

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/Opulentia-Trading/Arbitrage/contracts/uniswapV2Pair"
	"github.com/Opulentia-Trading/Arbitrage/platformListener"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethMath "github.com/ethereum/go-ethereum/common/math"
)

type UniswapV2Handler struct {
	*platformListener.EthHandler
}

type UniswapV2PairInfo struct {
	Symbol      string
	PairAddress string
	Token0      *platformListener.ERC20Token
	Token1      *platformListener.ERC20Token
	Reserve0    *big.Int
	Reserve1    *big.Int
}

func NewUniswapV2Handler() *UniswapV2Handler {
	infuraUrl := fmt.Sprintf("https://mainnet.infura.io/v3/%v", os.Getenv("INFURA_PROJECT_ID"))

	return &UniswapV2Handler{platformListener.NewEthHandler(infuraUrl)}
}

// Implements the ListenerHandler interface
func (u *UniswapV2Handler) TestConnection() {
	result := fmt.Sprintf("Latest block: %v", u.GetLatestBlockNumber())
	log.Println(result)
}

// Implements the ListenerHandler interface
func (u *UniswapV2Handler) FetchTickerPriceAll() []platformListener.TickerPrice {
	var result []platformListener.TickerPrice
	for _, pair := range uniswapV2PairMap {
		ticker := u.getPairPrice(pair)
		result = append(result, ticker)
	}

	return result
}

// Implements the ListenerHandler interface
func (u *UniswapV2Handler) FetchTickerPrice(asset1 string, asset2 string) platformListener.TickerPrice {
	symbol := fmt.Sprintf("%v/%v", asset1, asset2)
	pair, pairFound := uniswapV2PairMap[symbol]
	if !pairFound {
		log.Fatal("Unknown UniswapV2 symbol: ", symbol)
	}

	return u.getPairPrice(pair)
}

func (u *UniswapV2Handler) FetchPairInfo(asset1 string, asset2 string) *UniswapV2PairInfo {
	symbol := fmt.Sprintf("%v/%v", asset1, asset2)
	pair, pairFound := uniswapV2PairMap[symbol]
	if !pairFound {
		log.Fatal("Unknown UniswapV2 symbol: ", symbol)
	}

	instance := u.getPairInstance(pair.pairAddress)
	reserves, err := instance.GetReserves(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	return &UniswapV2PairInfo{
		Symbol:      symbol,
		PairAddress: pair.pairAddress,
		Token0:      pair.token0,
		Token1:      pair.token1,
		Reserve0:    reserves.Reserve0,
		Reserve1:    reserves.Reserve1,
	}
}

func (u *UniswapV2Handler) getPairInstance(address string) *uniswapV2Pair.UniswapV2Pair {
	// TODO: Implement caching of pair instances
	pairAddress := common.HexToAddress(address)
	instance, err := uniswapV2Pair.NewUniswapV2Pair(pairAddress, u.Client)
	if err != nil {
		log.Fatal(err)
	}

	return instance
}

// Returns the current mid price of a pair
func (u *UniswapV2Handler) getPairPrice(pair *pairWrapper) platformListener.TickerPrice {
	instance := u.getPairInstance(pair.pairAddress)
	reserves, err := instance.GetReserves(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	token0Price := new(big.Rat).SetFrac(reserves.Reserve1, reserves.Reserve0)

	// Determine the priceScalar based on the number of decimals in the
	// base and quote tokens. Computed using the formula below:
	// priceScalar = (10 ^ baseDecimals) / (10 ^ quoteDecimals)
	diffDecimals := pair.token0.Decimals - pair.token1.Decimals
	priceScalar := new(big.Rat)

	if diffDecimals >= 0 {
		priceScalar.SetInt(gethMath.BigPow(10, diffDecimals))
		token0Price.Mul(token0Price, priceScalar)
	} else {
		// Negative exponentiation is not supported by the big package
		// In this case, divide token0Price by priceScalar
		priceScalar.SetInt(gethMath.BigPow(10, -diffDecimals))
		token0Price.Quo(token0Price, priceScalar)
	}

	return platformListener.TickerPrice{
		Symbol: pair.token0.Symbol + pair.token1.Symbol,
		Price:  token0Price.FloatString(int(pair.token1.Decimals)),
	}
}
