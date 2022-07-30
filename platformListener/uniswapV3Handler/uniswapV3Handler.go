package uniswapV3Handler

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/Opulentia-Trading/Arbitrage/contracts/uniswapV3Pool"
	"github.com/Opulentia-Trading/Arbitrage/platformListener"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethMath "github.com/ethereum/go-ethereum/common/math"
)

var (
	q192 = gethMath.BigPow(2, 192)
)

type UniswapV3Handler struct {
	*platformListener.EthHandler
}

func NewUniswapV3Handler() *UniswapV3Handler {
	infuraUrl := fmt.Sprintf("https://mainnet.infura.io/v3/%v", os.Getenv("INFURA_PROJECT_ID"))

	return &UniswapV3Handler{platformListener.NewEthHandler(infuraUrl)}
}

// Implements the ListenerHandler interface
func (u *UniswapV3Handler) TestConnection() {
	result := fmt.Sprintf("Latest block: %v", u.GetLatestBlockNumber())
	log.Println(result)
}

// Implements the ListenerHandler interface
func (u *UniswapV3Handler) FetchTickerPriceAll() []platformListener.TickerPrice {
	var result []platformListener.TickerPrice

	for _, pool := range uniswapV3PoolMap {
		ticker := u.getPoolPrice(pool)
		result = append(result, ticker)
	}

	return result
}

// Implements the ListenerHandler interface
func (u *UniswapV3Handler) FetchTickerPrice(asset1 string, asset2 string) platformListener.TickerPrice {
	symbol := fmt.Sprintf("%v/%v", asset1, asset2)
	pool, poolFound := uniswapV3PoolMap[symbol]
	if !poolFound {
		log.Fatal("Unknown UniswapV3 symbol: ", symbol)
	}

	return u.getPoolPrice(pool)
}

func (u *UniswapV3Handler) getPoolInstance(address string) *uniswapV3Pool.UniswapV3Pool {
	// TODO: Implement caching of pool instances
	poolAddress := common.HexToAddress(address)
	instance, err := uniswapV3Pool.NewUniswapV3Pool(poolAddress, u.Client)
	if err != nil {
		log.Fatal(err)
	}

	return instance
}

// Returns the current mid price of a pool
func (u *UniswapV3Handler) getPoolPrice(pool *poolWrapper) platformListener.TickerPrice {
	instance := u.getPoolInstance(pool.poolAddress)
	poolState, err := instance.Slot0(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	// https://docs.uniswap.org/sdk/guides/fetching-prices
	// Convert SqrtPriceX96 to token0 price using the formula below:
	// price = SqrtPriceX96 ** 2 / 2 ** 192
	priceX96 := new(big.Int)
	priceX96.Mul(poolState.SqrtPriceX96, poolState.SqrtPriceX96)
	token0Price := new(big.Rat).SetFrac(priceX96, q192)

	// Determine the priceScalar based on the number of decimals in the
	// base and quote tokens. Computed using the formula below:
	// priceScalar = (10 ^ baseDecimals) / (10 ^ quoteDecimals)
	diffDecimals := pool.token0.decimals - pool.token1.decimals
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
		Symbol: pool.token0.symbol + pool.token1.symbol,
		Price:  token0Price.FloatString(int(pool.token1.decimals)),
	}
}
