package uniswapV3Handler

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/contracts/uniswapV3Pool"
	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethMath "github.com/ethereum/go-ethereum/common/math"
)

const PlatformName = "uniswap_v3"

var q192 = gethMath.BigPow(2, 192)

// Implements the Platform interface
type UniswapV3Handler struct {
	*ethHandler.EthHandler
}

func NewUniswapV3Handler() (*UniswapV3Handler, error) {
	exchangeInfo := models.Exchange{
		Type: models.Decentralized,
		Name: PlatformName,
	}

	network, err := ethHandler.GetEvmNetwork("ethereum_mainnet")
	if err != nil {
		panic(err)
	}

	if !isSupportedNetwork(network) {
		err := fmt.Errorf("unsupported network: %v", network)
		panic(err)
	}

	provider, err := ethHandler.GetEvmProvider("infura")
	if err != nil {
		panic(err)
	}

	providerProtocol := ethHandler.Https
	ethHandlerInst, err := ethHandler.NewEthHandler(network, provider, providerProtocol, &exchangeInfo)
	if err != nil {
		panic(err)
	}

	return &UniswapV3Handler{ethHandlerInst}, nil
}

func isSupportedNetwork(network *ethHandler.EvmNetwork) bool {
	supported := false

	switch network.ChainId {
	case 1, 5:
		supported = true // Ethereum
	case 137, 80001:
		supported = true // Polygon
	case 10, 420:
		supported = true // Optimism
	case 42161, 421613:
		supported = true // Arbitrum
	case 42220, 44787:
		supported = true // Celo
	}

	return supported
}

func (h *UniswapV3Handler) GetExchangeInfo() *models.Exchange {
	return h.ExchangeInfo
}

func (h *UniswapV3Handler) TestConnection() (string, error) {
	return h.EthHandler.TestConnection()
}

func (h *UniswapV3Handler) FetchTickerInfoAll() ([]models.TickerInfo, error) {
	var result []models.TickerInfo

	for _, pool := range poolsMap {
		if pool.ChainId == h.Network.ChainId {
			ticker, err := h.getPoolPrice(pool)
			if err != nil {
				panic(err)
			}
			result = append(result, ticker)
		}
	}

	return result, nil
}

func (h *UniswapV3Handler) FetchTickerInfo(base string, quote string) (models.TickerInfo, error) {
	// TODO: Handle other pool fee tiers
	// Default to the 0.3% fee tier
	var poolFee uint = 3000

	pool, err := GetPool(h.Network.ChainId, poolFee, base, quote)
	if err != nil {
		panic(err)
	}

	return h.getPoolPrice(pool)
}

// Returns an instance for interacting with the IUniswapV3Pool smart contract
func (h *UniswapV3Handler) getPoolInstance(address string) (*uniswapV3Pool.UniswapV3Pool, error) {
	poolAddress := common.HexToAddress(address)
	instance, err := uniswapV3Pool.NewUniswapV3Pool(poolAddress, h.Client)
	if err != nil {
		panic(err)
	}

	return instance, nil
}

// Returns the current mid price of a pool
func (h *UniswapV3Handler) getPoolPrice(pool *PoolWrapper) (models.TickerInfo, error) {
	instance, err := h.getPoolInstance(pool.PoolAddress)
	if err != nil {
		panic(err)
	}

	token0, err := ethHandler.GetToken(pool.ChainId, pool.Token0Symbol)
	if err != nil {
		panic(err)
	}

	token1, err := ethHandler.GetToken(pool.ChainId, pool.Token1Symbol)
	if err != nil {
		panic(err)
	}

	poolState, err := instance.Slot0(&bind.CallOpts{})
	if err != nil {
		panic(err)
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
	diffDecimals := int64(token0.Decimals) - int64(token1.Decimals)
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

	result := models.TickerInfo{
		Symbol:         pool.Symbol(),
		Base:           pool.Token0Symbol,
		Quote:          pool.Token1Symbol,
		Price:          token0Price.FloatString(int(token1.Decimals)),
		MakerComission: pool.FeeString(),
		TakerComission: pool.FeeString(),
		Timestamp:      time.Now(),
	}

	return result, nil
}

func (h *UniswapV3Handler) ExecuteOrder(order models.Order) error {
	fmt.Printf("Executing %v/%v %v order\n", order.Base, order.Quote, order.Action.String())
	return nil
}

func (h *UniswapV3Handler) String() string {
	return h.ExchangeInfo.Name
}
