package platformListener

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthHandler struct {
	ProviderUrl string
	Client      *ethclient.Client
}

func NewEthHandler(providerUrl string) *EthHandler {
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		log.Fatal(err)
	}

	return &EthHandler{ProviderUrl: providerUrl, Client: client}
}

func (e *EthHandler) GetLatestBlockNumber() string {
	header, err := e.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return header.Number.String()
}
