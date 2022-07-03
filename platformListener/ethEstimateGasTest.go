package platformListener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Opulentia-Trading/Arbitrage/util"
	gethMath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
Eth Gas Reference:
https://ethereum.org/en/developers/docs/gas/
https://www.blocknative.com/blog/eip-1559-fees

Post London Upgrade (EIP-1559), total transaction fee is calculated using:
Gas units (Gas limit) * (Base fee + Tip)

Gas limit
---------
The gas limit depends on the type of transaction. For a simple ETH transfer it is usually 21,000 units.
For more complex smart contract interactions, it will be higher.

Base Fee
--------
The base fee is burnt and not shared with miners. The base fee for the next pending block
is determined automatically by by the blockchain so we don't do any estimations here.

Priority Fee (Tip)
------------------
This fee is sent directly to miners as a tip. We have to determine a suitable tip based on previous blocks and
pending transactions in the mempool. The tip should be high enough so that our transaction is included in the next block, but
not too high as this will reduce profits.

Max Fee
-------
There is another param which determines the absolute maximum price we are willing to pay for each gas step.
Any gas not used in the transaction is returned as follows:
refund = max fee - (base fee + tip)

Imagine we submit a transaction for block i but it is not included within this block. For the next block i+1,
the base fee has increased. With a max fee, we can handle this scenario and possibly include the transaction in
block i+1.
*/

var (
	gweiToWei = gethMath.BigPow(10, 9)
)

type gasEstimate struct {
	Source      string
	Latency     int64 // ms
	BlockNumber *big.Int
	BaseFee     *big.Int // Wei
	PriorityFee *big.Int // Wei
	MaxFee      *big.Int // Wei
}

type blocknativeResponse struct {
	System             string `json:"system"`
	Network            string `json:"network"`
	Unit               string `json:"unit"`
	MaxPrice           int    `json:"maxPrice"`
	CurrentBlockNumber int    `json:"currentBlockNumber"`
	MsSinceLastBlock   int    `json:"msSinceLastBlock"`
	BlockPrices        []struct {
		BlockNumber               int64   `json:"blockNumber"`
		EstimatedTransactionCount int     `json:"estimatedTransactionCount"`
		BaseFeePerGas             float64 `json:"baseFeePerGas"`
		EstimatedPrices           []struct {
			Confidence           int     `json:"confidence"`
			Price                int     `json:"price"`
			MaxPriorityFeePerGas float64 `json:"maxPriorityFeePerGas"`
			MaxFeePerGas         float64 `json:"maxFeePerGas"`
		} `json:"estimatedPrices"`
	} `json:"blockPrices"`
	EstimatedBaseFees []struct {
		Pending1 []blocknativeBaseFeeEstimate `json:"pending+1,omitempty"`
		Pending2 []blocknativeBaseFeeEstimate `json:"pending+2,omitempty"`
		Pending3 []blocknativeBaseFeeEstimate `json:"pending+3,omitempty"`
		Pending4 []blocknativeBaseFeeEstimate `json:"pending+4,omitempty"`
		Pending5 []blocknativeBaseFeeEstimate `json:"pending+5,omitempty"`
	} `json:"estimatedBaseFees"`
}

type blocknativeBaseFeeEstimate struct {
	Confidence int     `json:"confidence"`
	BaseFee    float64 `json:"baseFee"`
}

type etherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		LastBlock       string `json:"LastBlock"`
		SafeGasPrice    string `json:"SafeGasPrice"`
		ProposeGasPrice string `json:"ProposeGasPrice"`
		FastGasPrice    string `json:"FastGasPrice"`
		SuggestBaseFee  string `json:"suggestBaseFee"`
		GasUsedRatio    string `json:"gasUsedRatio"`
	} `json:"result"`
}

func blocknativeGasTest() *gasEstimate {
	// Blocknative Gas Platform
	// https://www.blocknative.com/blog/introducing-gas-platform
	// https://www.blocknative.com/blog/comparing-eth-gas-estimators
	// Uses a quantile regression model to estimate gas prices based on the mempool and previous blocks
	// Provides gas estimates for different cofidence levels (99%, 95%, 90%, 80%, and 70%)
	url := "https://api.blocknative.com/gasprices/blockprices"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", os.Getenv("BLOCKNATIVE_API_KEY"))

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var apiResponse blocknativeResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&apiResponse)
	if err != nil {
		log.Fatalln(err)
	}
	elapsed := time.Since(start)

	result := gasEstimate{
		Source:      "Blocknative",
		Latency:     elapsed.Milliseconds(),
		BlockNumber: big.NewInt(apiResponse.BlockPrices[0].BlockNumber),
		BaseFee:     convGweiToWei(apiResponse.BlockPrices[0].BaseFeePerGas),
	}

	// Use values from the 95% confidence level
	for _, estimatedPrice := range apiResponse.BlockPrices[0].EstimatedPrices {
		if estimatedPrice.Confidence == 95 {
			result.PriorityFee = convGweiToWei(estimatedPrice.MaxPriorityFeePerGas)
			result.MaxFee = convGweiToWei(estimatedPrice.MaxFeePerGas)
			break
		}
	}

	return &result

}

func etherscanGasTest() *gasEstimate {
	// Etherscan Gas API
	// Docs don't specify the estimation technique
	// However, Blocknative thinks they use a time based approach
	url := fmt.Sprintf("https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=%v", os.Getenv("ETHERSCAN_API_KEY"))
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var apiResponse etherscanResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&apiResponse)
	if err != nil {
		log.Fatalln(err)
	}
	elapsed := time.Since(start)

	blockNum, ok := new(big.Int).SetString(apiResponse.Result.LastBlock, 10)
	if !ok {
		log.Fatalln("Failed to parse block number: ", apiResponse.Result.LastBlock)
	}
	blockNum.Add(blockNum, big.NewInt(1)) // increment for pending block

	baseFee := convGweiToWei(apiResponse.Result.SuggestBaseFee)
	fastGasPrice := convGweiToWei(apiResponse.Result.FastGasPrice)
	priorityFee := new(big.Int).Sub(fastGasPrice, baseFee)

	return &gasEstimate{
		Source:      "Etherscan",
		Latency:     elapsed.Milliseconds(),
		BlockNumber: blockNum,
		BaseFee:     baseFee,
		PriorityFee: priorityFee,
		MaxFee:      nil, // Not provided by API
	}
}

func gethGasTest() *gasEstimate {
	// Estimate using Go Ethereum (Geth) packages
	providerUrl := fmt.Sprintf("https://mainnet.infura.io/v3/%v", os.Getenv("INFURA_PROJECT_ID"))
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	start := time.Now()
	pendingBlockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(-1))
	if err != nil {
		log.Fatal(err)
	}

	// SuggestGasTipCap retrieves the cheapest 3 transactions from the past X blocks (X = 20 for full nodes; X = 2 for light clients),
	// and uses the 60th percentile as the suggestion for the priority fee.

	// Note: The predictions from SuggestGasTipCap seem to be underpriced for Arbitrage uses
	// A custom approach could use the eth_feeHistory API and mempool information from BloxRoute to make better predictions
	// Refer to https://docs.alchemy.com/alchemy/guides/eip-1559/gas-estimator for a basic estimator using the eth_feeHistory API
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	// Simple heuristic for estimating the max fee
	// Max Fee = (2 * Base Fee) + Max Priority Fee
	// Ensures the fee will be competitive for six consecutive blocks,
	// assuming the base fee for each subsequent block increases by the maximum of 12.5%
	maxFee := new(big.Int).Mul(pendingBlockHeader.BaseFee, big.NewInt(2))
	maxFee.Add(maxFee, gasTipCap)

	return &gasEstimate{
		Source:      "Go Ethereum",
		Latency:     elapsed.Milliseconds(),
		BlockNumber: pendingBlockHeader.Number,
		BaseFee:     pendingBlockHeader.BaseFee,
		PriorityFee: gasTipCap,
		MaxFee:      maxFee,
	}
}

func convGweiToWei(gwei interface{}) *big.Int {
	bigGwei := new(big.Float)

	switch val := gwei.(type) {
	case float64:
		bigGwei.SetFloat64(val)
	case string:
		_, ok := bigGwei.SetString(val)
		if !ok {
			log.Fatalln("Failed to parse value: ", val)
		}
	default:
		log.Fatalln("Unsupported type")
	}

	scalar := new(big.Float).SetInt(gweiToWei)
	bigGwei.Mul(bigGwei, scalar)
	result, _ := bigGwei.Int(nil)
	return result
}

func RunEthGasTests() {
	blocknative := blocknativeGasTest()
	etherscan := etherscanGasTest()
	geth := gethGasTest()

	log.Println(util.PrettyPrint(blocknative))
	log.Println(util.PrettyPrint(etherscan))
	log.Println(util.PrettyPrint(geth))
}
