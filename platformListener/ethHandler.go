package platformListener

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type ethHandler struct {
	providerUrl string
	client      *ethclient.Client
}

func newEthHandler(providerUrl string) *ethHandler {
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		log.Fatal(err)
	}

	return &ethHandler{providerUrl: providerUrl, client: client}
}

func (e *ethHandler) getLatestBlockNumber() string {
	header, err := e.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return header.Number.String()
}
