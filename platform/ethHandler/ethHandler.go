package ethHandler

import (
	"context"
	"fmt"

	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthHandler struct {
	Network          *EvmNetwork
	Provider         EvmProvider
	ProviderProtocol ProviderProtocol
	ExchangeInfo     *models.Exchange
	Client           *ethclient.Client
}

func NewEthHandler(
	network *EvmNetwork,
	provider EvmProvider,
	providerProtocol ProviderProtocol,
	exchangeInfo *models.Exchange,
) (*EthHandler, error) {
	rpcEndpoints, err := provider.GetRpcEndpoints(network.ChainId, providerProtocol)
	if err != nil {
		panic(err)
	}

	client, err := ethclient.Dial(rpcEndpoints[0])
	if err != nil {
		panic(err)
	}

	ethHandler := &EthHandler{
		Network:          network,
		Provider:         provider,
		ProviderProtocol: providerProtocol,
		ExchangeInfo:     exchangeInfo,
		Client:           client,
	}

	return ethHandler, nil
}

func (e *EthHandler) GetLatestBlockNumber() (string, error) {
	header, err := e.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return header.Number.String(), nil
}

func (e *EthHandler) TestConnection() (string, error) {
	latestBlock, err := e.GetLatestBlockNumber()
	if err != nil {
		panic(err)
	}

	output := fmt.Sprintf("Network: %v\nProvider: %v\nProvider protocol: %v\nLatest block: %v",
		e.Network,
		e.Provider,
		e.ProviderProtocol,
		latestBlock)

	return output, nil
}
