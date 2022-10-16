package ethHandler

import (
	"fmt"
	"strings"
)

type ChainId uint

type NativeCurrency struct {
	Name     string
	Symbol   string
	Decimals uint8
}

type EvmNetwork struct {
	Name           string
	IsMainnet      bool
	ChainId        ChainId
	NativeCurrency *NativeCurrency
}

// Network Listing: https://chainid.network/chains.json
// TODO: Parse from json or config file
var evmNetworkMap = map[string]*EvmNetwork{
	"ethereum_mainnet": {
		Name:      "ethereum_mainnet",
		IsMainnet: true,
		ChainId:   ChainId(1),
		NativeCurrency: &NativeCurrency{
			Name:     "Ether",
			Symbol:   "ETH",
			Decimals: 18,
		},
	},
	"ethereum_goerli": {
		Name:      "ethereum_goerli",
		IsMainnet: false,
		ChainId:   ChainId(5),
		NativeCurrency: &NativeCurrency{
			Name:     "Goerli Ether",
			Symbol:   "ETH",
			Decimals: 18,
		},
	},
}

func GetEvmNetwork(networkName string) (*EvmNetwork, error) {
	networkName = strings.ToLower(networkName)
	network, networkFound := evmNetworkMap[networkName]
	if !networkFound {
		return nil, fmt.Errorf("unknown network: %v", networkName)
	}

	return network, nil
}

func (e *EvmNetwork) String() string {
	out := fmt.Sprintf("%v(chainId=%v)", e.Name, e.ChainId)
	return out
}
