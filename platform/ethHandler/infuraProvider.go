package ethHandler

import (
	"fmt"
	"os"
	"sync"
)

const InfuraProviderName = "infura"

// Implements EvmProvider
type InfuraProvider struct {
	endpointsMap map[string][]string
}

var infuraProviderOnce sync.Once
var infuraProviderInst *InfuraProvider

func genEndpointsMap() map[string][]string {
	endpointsMap := map[string][]string{
		genEndpointsMapKey(1, Https):      {"https://mainnet.infura.io/v3"},
		genEndpointsMapKey(1, WebSockets): {"wss://mainnet.infura.io/ws/v3"},
		genEndpointsMapKey(5, Https):      {"https://goerli.infura.io/v3"},
		genEndpointsMapKey(5, WebSockets): {"wss://goerli.infura.io/ws/v3"},
	}

	for key, rpcUrls := range endpointsMap {
		endpointsMap[key] = formatRpcUrls(rpcUrls)
	}

	return endpointsMap
}

func genEndpointsMapKey(chainId ChainId, protocol ProviderProtocol) string {
	return fmt.Sprintf("%v|%v", chainId, protocol)
}

func formatRpcUrls(urls []string) []string {
	result := make([]string, len(urls))
	for i, url := range urls {
		formattedUrl := url + "/" + os.Getenv("INFURA_PROJECT_ID")
		result[i] = formattedUrl
	}

	return result
}

func GetInfuraProvider() *InfuraProvider {
	infuraProviderOnce.Do(func() {
		endpointsMap := genEndpointsMap()
		infuraProviderInst = &InfuraProvider{endpointsMap: endpointsMap}
	})

	return infuraProviderInst
}

func (p *InfuraProvider) GetRpcEndpoints(chainId ChainId, protocol ProviderProtocol) ([]string, error) {
	key := genEndpointsMapKey(chainId, protocol)
	endpoints, found := p.endpointsMap[key]
	if !found {
		return nil, fmt.Errorf("unknown endpoints with chainId=%v protocol=%v", chainId, protocol)
	}

	return endpoints, nil
}

func (p *InfuraProvider) String() string {
	return InfuraProviderName
}
