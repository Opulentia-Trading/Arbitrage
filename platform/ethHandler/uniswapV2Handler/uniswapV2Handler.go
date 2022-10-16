package uniswapV2Handler

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/contracts/uniswapV2Pair"
	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethMath "github.com/ethereum/go-ethereum/common/math"
)

const PlatformName = "uniswap_v2"

// Implements the Platform interface
type UniswapV2Handler struct {
	*ethHandler.EthHandler
}

type PairReserves struct {
	Symbol             string
	ChainId            ethHandler.ChainId
	PairAddress        string
	Token0             *ethHandler.Token
	Token1             *ethHandler.Token
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast time.Time
}

func NewUniswapV2Handler() (*UniswapV2Handler, error) {
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

	return &UniswapV2Handler{ethHandlerInst}, nil
}

func isSupportedNetwork(network *ethHandler.EvmNetwork) bool {
	return (network.ChainId == 1 || network.ChainId == 5) // Ethereum
}

func (h *UniswapV2Handler) GetExchangeInfo() *models.Exchange {
	return h.ExchangeInfo
}

func (h *UniswapV2Handler) TestConnection() (string, error) {
	return h.EthHandler.TestConnection()
}

func (h *UniswapV2Handler) FetchTickerInfoAll() ([]models.TickerInfo, error) {
	var result []models.TickerInfo

	for _, pair := range pairsMap {
		if pair.ChainId == h.Network.ChainId {
			ticker, err := h.getPairPrice(pair)
			if err != nil {
				panic(err)
			}
			result = append(result, ticker)
		}
	}

	return result, nil
}

func (h *UniswapV2Handler) FetchTickerInfo(base string, quote string) (models.TickerInfo, error) {
	pair, err := GetPair(h.Network.ChainId, base, quote)
	if err != nil {
		panic(err)
	}

	return h.getPairPrice(pair)
}

// Returns an instance for interacting with the IUniswapV2Pair smart contract
func (h *UniswapV2Handler) getPairInstance(address string) (*uniswapV2Pair.UniswapV2Pair, error) {
	pairAddress := common.HexToAddress(address)
	instance, err := uniswapV2Pair.NewUniswapV2Pair(pairAddress, h.Client)
	if err != nil {
		panic(err)
	}

	return instance, nil
}

// Returns the current mid price of a pair
func (h *UniswapV2Handler) getPairPrice(pair *PairWrapper) (models.TickerInfo, error) {
	var result models.TickerInfo

	instance, err := h.getPairInstance(pair.PairAddress)
	if err != nil {
		panic(err)
	}

	token0, err := ethHandler.GetToken(pair.ChainId, pair.Token0Symbol)
	if err != nil {
		panic(err)
	}

	token1, err := ethHandler.GetToken(pair.ChainId, pair.Token1Symbol)
	if err != nil {
		panic(err)
	}

	reserves, err := instance.GetReserves(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	token0Price := new(big.Rat).SetFrac(reserves.Reserve1, reserves.Reserve0)

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

	result = models.TickerInfo{
		Symbol:         pair.Symbol(),
		Base:           pair.Token0Symbol,
		Quote:          pair.Token1Symbol,
		Price:          token0Price.FloatString(int(token1.Decimals)),
		MakerComission: "0.3",
		TakerComission: "0.3",
		Timestamp:      time.Now(),
	}

	return result, nil
}

func (h *UniswapV2Handler) FetchPairReserves(base string, quote string) (*PairReserves, error) {
	pair, err := GetPair(h.Network.ChainId, base, quote)
	if err != nil {
		panic(err)
	}

	instance, err := h.getPairInstance(pair.PairAddress)
	if err != nil {
		panic(err)
	}

	token0, err := ethHandler.GetToken(pair.ChainId, pair.Token0Symbol)
	if err != nil {
		panic(err)
	}

	token1, err := ethHandler.GetToken(pair.ChainId, pair.Token1Symbol)
	if err != nil {
		panic(err)
	}

	reserves, err := instance.GetReserves(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	result := &PairReserves{
		Symbol:             pair.Symbol(),
		ChainId:            pair.ChainId,
		PairAddress:        pair.PairAddress,
		Token0:             token0,
		Token1:             token1,
		Reserve0:           reserves.Reserve0,
		Reserve1:           reserves.Reserve1,
		BlockTimestampLast: time.Unix(int64(reserves.BlockTimestampLast), 0),
	}

	return result, nil
}

func (h *UniswapV2Handler) ExecuteOrder(order models.Order) error {
	fmt.Printf("Executing %v/%v %v order\n", order.Base, order.Quote, order.Action.String())
	return nil
}

func (h *UniswapV2Handler) String() string {
	return h.ExchangeInfo.Name
}
