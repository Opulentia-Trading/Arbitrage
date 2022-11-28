package uniswapV2Handler

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/contracts/uniswapV2Pair"
	"github.com/Opulentia-Trading/Arbitrage/contracts/uniswapV2Router02"
	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethMath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	PlatformName      = "uniswap_v2"
	swapFee           = "0.3" // in percent
	txMineWaitTimeout = 5 * time.Minute
)

var router02Address = common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")

// Implements the Platform interface
type UniswapV2Handler struct {
	*ethHandler.EthHandler
	SwapNativeETH bool // use native ETH as the input/output of a swap
	SendSwapTx    bool // broadcast swap tx on blockchain
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

	network, err := ethHandler.GetEvmNetwork("ethereum_goerli")
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

	return &UniswapV2Handler{
		EthHandler:    ethHandlerInst,
		SwapNativeETH: false,
		SendSwapTx:    true,
	}, nil
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

	result := models.TickerInfo{
		Symbol:         pair.Symbol(),
		Base:           pair.Token0Symbol,
		Quote:          pair.Token1Symbol,
		Price:          token0Price.FloatString(int(token1.Decimals)),
		MakerComission: swapFee,
		TakerComission: swapFee,
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

// Returns an instance for interacting with the IUniswapV2Router02 smart contract
func (h *UniswapV2Handler) getRouter02Instance() (*uniswapV2Router02.UniswapV2Router02, error) {
	instance, err := uniswapV2Router02.NewUniswapV2Router02(router02Address, h.Client)
	if err != nil {
		panic(err)
	}

	return instance, nil
}

func (h *UniswapV2Handler) getOrderPath(
	base *ethHandler.Token,
	quote *ethHandler.Token,
	action models.Action,
) ([]common.Address, *ethHandler.Token, error) {
	baseAddress := base.AddressForGeth()
	quoteAddress := quote.AddressForGeth()
	var path []common.Address
	var inputToken *ethHandler.Token

	switch action {
	case models.BuyLongSpot:
		path = []common.Address{quoteAddress, baseAddress}
		inputToken = quote
	case models.SellLongSpot:
		path = []common.Address{baseAddress, quoteAddress}
		inputToken = base
	default:
		return nil, nil, fmt.Errorf("unsupported action %v", action)
	}

	return path, inputToken, nil
}

func (h *UniswapV2Handler) approveToken(wallet *ethHandler.Wallet, token *ethHandler.Token) error {
	tokenHandler, err := ethHandler.NewERC20Handler(h.EthHandler, token)
	if err != nil {
		panic(err)
	}

	err = tokenHandler.MaxApprove(wallet, router02Address, true)
	if err != nil {
		panic(err)
	}

	return nil
}

func (h *UniswapV2Handler) ExecuteOrder(order models.Order) error {
	wallet, err := ethHandler.GetWallet(os.Getenv("WALLET_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	chainId := big.NewInt(int64(h.Network.ChainId))
	auth, err := bind.NewKeyedTransactorWithChainID(wallet.PrivateKey, chainId)
	if err != nil {
		panic(err)
	}

	routerInstance, err := h.getRouter02Instance()
	if err != nil {
		panic(err)
	}

	baseSymbol, quoteSymbol := NormalizePairTokens(order.Base, order.Quote)
	baseToken, err := ethHandler.GetToken(h.Network.ChainId, baseSymbol)
	if err != nil {
		panic(err)
	}

	quoteToken, err := ethHandler.GetToken(h.Network.ChainId, quoteSymbol)
	if err != nil {
		panic(err)
	}

	wethToken, err := ethHandler.GetToken(h.Network.ChainId, "WETH")
	if err != nil {
		panic(err)
	}

	wethAddress := wethToken.AddressForGeth()
	path, inputToken, err := h.getOrderPath(baseToken, quoteToken, order.Action)
	if err != nil {
		panic(err)
	}

	if !h.SwapNativeETH || path[0] != wethAddress {
		// TODO: Pre-approve tokens on init
		err := h.approveToken(wallet, inputToken)
		if err != nil {
			panic(err)
		}
	}

	// TODO: Maybe keep track of nonce locally
	nonce, err := h.Client.PendingNonceAt(context.Background(), wallet.Address)
	if err != nil {
		panic(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = nil
	auth.NoSend = !h.SendSwapTx

	// TODO: Get gas estimates from the gasEstimator module
	auth.GasPrice = nil
	auth.GasFeeCap = nil
	auth.GasTipCap = nil
	auth.GasLimit = uint64(0) // in units (300000 should be a good upper bound)

	deadline := big.NewInt(time.Now().Add(order.Deadline).Unix())
	var tx *types.Transaction = nil

	if h.SwapNativeETH && path[0] == wethAddress {
		auth.Value = order.LiqPoolAmountIn
		tx, err = routerInstance.SwapExactETHForTokens(
			auth,
			order.LiqPoolAmountOut,
			path,
			wallet.Address,
			deadline)

		if err != nil {
			panic(err)
		}
	} else if h.SwapNativeETH && path[len(path)-1] == wethAddress {
		tx, err = routerInstance.SwapExactTokensForETH(
			auth,
			order.LiqPoolAmountIn,
			order.LiqPoolAmountOut,
			path,
			wallet.Address,
			deadline)

		if err != nil {
			panic(err)
		}
	} else {
		tx, err = routerInstance.SwapExactTokensForTokens(
			auth,
			order.LiqPoolAmountIn,
			order.LiqPoolAmountOut,
			path,
			wallet.Address,
			deadline)

		if err != nil {
			panic(err)
		}
	}

	if tx == nil {
		panic("failed to prepare transaction")
	}

	fmt.Printf("\n[[ %v/%v %v tx ]]\n", order.Base, order.Quote, order.Action.String())
	fmt.Printf("tx hash: %s\n", tx.Hash())
	fmt.Printf("gas priority fee: %v\n", tx.GasTipCap())
	fmt.Printf("gas max fee: %v\n", tx.GasFeeCap())
	fmt.Printf("gas limit: %v\n", tx.Gas())
	if auth.NoSend {
		fmt.Println("Note: transaction not sent on blockchain")
		return nil
	}

	_, err = h.WaitTxMined(tx, wallet.Address, txMineWaitTimeout)
	if err != nil {
		panic(err)
	}

	return nil
}

func (h *UniswapV2Handler) String() string {
	return h.ExchangeInfo.Name
}
