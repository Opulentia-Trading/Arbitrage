package ethHandler

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func (e *EthHandler) WaitTxMined(tx *types.Transaction, fromAddress common.Address, waitTimeout time.Duration) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	fmt.Printf("waiting for tx to be mined (waitTimeout=%v) ...\n", waitTimeout)
	txReceipt, err := bind.WaitMined(ctx, e.Client, tx)
	if err != nil {
		panic(err)
	}

	txSuccess := txReceipt.Status == types.ReceiptStatusSuccessful
	fmt.Println("\n[mined tx receipt]")
	fmt.Println("status success: ", txSuccess)
	fmt.Println("block number: ", txReceipt.BlockNumber)
	fmt.Println("block hash: ", txReceipt.BlockHash)
	fmt.Println("tx index: ", txReceipt.TransactionIndex)
	fmt.Println("gas used: ", txReceipt.GasUsed)
	fmt.Println("type: ", txReceipt.Type)

	if !txSuccess {
		err := e.FailedTxError(tx, fromAddress, txReceipt.BlockNumber)
		if err != nil {
			panic(err)
		}

		err = fmt.Errorf("txHash=%v failed", tx.Hash())
		panic(err)
	}

	return txReceipt, nil
}

// The error from a failed transaction is not included in the tx receipt or logs.
// The eth_call RPC method executes a message call directly in the VM of a node without creating a blockchain transaction.
// Using this, we can replay the failed transaction locally on a node to retrieve the error.
// Note: eth_call does not consume gas.
func (e *EthHandler) FailedTxError(tx *types.Transaction, fromAddress common.Address, blockNumber *big.Int) error {
	msg := ethereum.CallMsg{
		From:      fromAddress,
		To:        tx.To(),
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice(),
		GasFeeCap: tx.GasFeeCap(),
		GasTipCap: tx.GasTipCap(),
		Value:     tx.Value(),
		Data:      tx.Data(),
	}

	_, err := e.Client.CallContract(context.Background(), msg, blockNumber)
	return err
}
