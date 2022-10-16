package ethHandler

import "fmt"

type TokenType uint

const (
	ERC20 TokenType = iota
	ERC721
	ERC1155
)

func (t TokenType) String() string {
	return [...]string{
		"ERC20",
		"ERC721",
		"ERC1155"}[t]
}

type Token struct {
	ChainId  ChainId
	Type     TokenType
	Address  string
	Name     string
	Symbol   string
	Decimals uint8
}

// TODO: Parse from json or config file
var tokensMap = map[string]*Token{
	genTokensMapKey(1, "WETH"): {
		ChainId:  1,
		Type:     ERC20,
		Address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Name:     "Wrapped Ether",
		Symbol:   "WETH",
		Decimals: 18,
	},
	genTokensMapKey(1, "USDC"): {
		ChainId:  1,
		Type:     ERC20,
		Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Name:     "USD Coin",
		Symbol:   "USDC",
		Decimals: 6,
	},
	genTokensMapKey(1, "WBTC"): {
		ChainId:  1,
		Type:     ERC20,
		Address:  "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
		Name:     "Wrapped BTC",
		Symbol:   "WBTC",
		Decimals: 8,
	},
	genTokensMapKey(1, "APE"): {
		ChainId:  1,
		Type:     ERC20,
		Address:  "0x4d224452801ACEd8B2F0aebE155379bb5D594381",
		Name:     "ApeCoin",
		Symbol:   "APE",
		Decimals: 18,
	},
	genTokensMapKey(1, "LINK"): {
		ChainId:  1,
		Type:     ERC20,
		Address:  "0x514910771AF9Ca656af840dff83E8264EcF986CA",
		Name:     "ChainLink Token",
		Symbol:   "LINK",
		Decimals: 18,
	},
	genTokensMapKey(5, "WETH"): {
		ChainId:  5,
		Type:     ERC20,
		Address:  "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6",
		Name:     "Wrapped Ether",
		Symbol:   "WETH",
		Decimals: 18,
	},
	genTokensMapKey(5, "USDC"): {
		ChainId:  5,
		Type:     ERC20,
		Address:  "0x07865c6E87B9F70255377e024ace6630C1Eaa37F",
		Name:     "USD Coin",
		Symbol:   "USDC",
		Decimals: 6,
	},
	genTokensMapKey(5, "LINK"): {
		ChainId:  5,
		Type:     ERC20,
		Address:  "0x326C977E6efc84E512bB9C30f76E30c160eD06FB",
		Name:     "ChainLink Token",
		Symbol:   "LINK",
		Decimals: 18,
	},
}

func genTokensMapKey(chainId ChainId, symbol string) string {
	return fmt.Sprintf("%v|%v", chainId, symbol)
}

func GetToken(chainId ChainId, symbol string) (*Token, error) {
	key := genTokensMapKey(chainId, symbol)
	token, tokenFound := tokensMap[key]
	if !tokenFound {
		return nil, fmt.Errorf("unknown token with chainId=%v symbol=%v", chainId, symbol)
	}

	return token, nil
}

func (t *Token) String() string {
	out := fmt.Sprintf("%v(chainId=%v)", t.Symbol, t.ChainId)
	return out
}
