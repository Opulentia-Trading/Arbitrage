package ethHandler

import (
	"fmt"
	"strings"
)

type ProviderProtocol uint

const (
	Https ProviderProtocol = iota
	WebSockets
)

func (p ProviderProtocol) String() string {
	return [...]string{
		"Https",
		"WebSockets"}[p]
}

type EvmProvider interface {
	GetRpcEndpoints(chainId ChainId, protocol ProviderProtocol) ([]string, error)
	String() string
}

func GetEvmProvider(providerName string) (EvmProvider, error) {
	providerName = strings.ToLower(providerName)

	switch providerName {
	case InfuraProviderName:
		return GetInfuraProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %v", providerName)
	}
}
