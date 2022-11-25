package uniswapV2Handler

import (
	"fmt"

	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler"
)

type PairWrapper struct {
	ChainId      ethHandler.ChainId
	PairAddress  string
	Token0Symbol string
	Token1Symbol string
}

// TODO: Parse from json or config file
var pairsMap = map[string]*PairWrapper{
	genPairsMapKey(1, "USDC", "WETH"): {
		ChainId:      1,
		PairAddress:  "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
		Token0Symbol: "USDC",
		Token1Symbol: "WETH",
	},
	genPairsMapKey(1, "LINK", "WETH"): {
		ChainId:      1,
		PairAddress:  "0xa2107FA5B38d9bbd2C461D6EDf11B11A50F6b974",
		Token0Symbol: "LINK",
		Token1Symbol: "WETH",
	},
	genPairsMapKey(5, "USDC", "WETH"): {
		ChainId:      5,
		PairAddress:  "0x647595535c370F6092C6daE9D05a7Ce9A8819F37",
		Token0Symbol: "USDC",
		Token1Symbol: "WETH",
	},
	genPairsMapKey(5, "LINK", "WETH"): {
		ChainId:      5,
		PairAddress:  "0x32bE40dC4Db907aCf18773bfC81F1bFFA92B77c2",
		Token0Symbol: "LINK",
		Token1Symbol: "WETH",
	},
}

func genPairsMapKey(chainId ethHandler.ChainId, base string, quote string) string {
	return fmt.Sprintf("%v|%v/%v", chainId, base, quote)
}

func NormalizePairTokens(base string, quote string) (string, string) {
	// Wrap native ETH
	if base == "ETH" {
		base = "WETH"
	}

	if quote == "ETH" {
		quote = "WETH"
	}

	return base, quote
}

func GetPair(chainId ethHandler.ChainId, base string, quote string) (*PairWrapper, error) {
	base, quote = NormalizePairTokens(base, quote)
	key := genPairsMapKey(chainId, base, quote)
	pair, pairFound := pairsMap[key]
	if !pairFound {
		return nil, fmt.Errorf("unknown pair with chainId=%v base=%v quote=%v", chainId, base, quote)
	}

	return pair, nil
}

func (p *PairWrapper) Symbol() string {
	return p.Token0Symbol + "/" + p.Token1Symbol
}

func (p *PairWrapper) String() string {
	out := fmt.Sprintf("%v(chainId=%v)", p.Symbol(), p.ChainId)
	return out
}
