package uniswapV3Handler

import (
	"fmt"
	"strconv"

	"github.com/Opulentia-Trading/Arbitrage/platform/ethHandler"
)

var poolFeeToPercent = 1e-4

type PoolWrapper struct {
	ChainId      ethHandler.ChainId
	PoolAddress  string
	Fee          uint // fee is in hundredths of a basis point (0.3% => 3000)
	Token0Symbol string
	Token1Symbol string
}

// TODO: Parse from json or config file
var poolsMap = map[string]*PoolWrapper{
	genPoolsMapKey(1, 3000, "USDC", "WETH"): {
		ChainId:      1,
		PoolAddress:  "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
		Fee:          3000,
		Token0Symbol: "USDC",
		Token1Symbol: "WETH",
	},
	genPoolsMapKey(1, 3000, "WBTC", "WETH"): {
		ChainId:      1,
		PoolAddress:  "0xCBCdF9626bC03E24f779434178A73a0B4bad62eD",
		Fee:          3000,
		Token0Symbol: "WBTC",
		Token1Symbol: "WETH",
	},
	genPoolsMapKey(1, 3000, "APE", "WETH"): {
		ChainId:      1,
		PoolAddress:  "0xAc4b3DacB91461209Ae9d41EC517c2B9Cb1B7DAF",
		Fee:          3000,
		Token0Symbol: "APE",
		Token1Symbol: "WETH",
	},
	genPoolsMapKey(1, 3000, "LINK", "WETH"): {
		ChainId:      1,
		PoolAddress:  "0xa6Cc3C2531FdaA6Ae1A3CA84c2855806728693e8",
		Fee:          3000,
		Token0Symbol: "LINK",
		Token1Symbol: "WETH",
	},
}

func genPoolsMapKey(chainId ethHandler.ChainId, fee uint, base string, quote string) string {
	return fmt.Sprintf("%v|%v|%v/%v", chainId, fee, base, quote)
}

func GetPool(chainId ethHandler.ChainId, fee uint, base string, quote string) (*PoolWrapper, error) {
	// Wrap native ETH
	if base == "ETH" {
		base = "WETH"
	}
	if quote == "ETH" {
		quote = "WETH"
	}

	key := genPoolsMapKey(chainId, fee, base, quote)
	pool, poolFound := poolsMap[key]
	if !poolFound {
		return nil, fmt.Errorf("unknown pool with chainId=%v fee=%v base=%v quote=%v", chainId, fee, base, quote)
	}

	return pool, nil
}

func (p *PoolWrapper) Symbol() string {
	return p.Token0Symbol + "/" + p.Token1Symbol
}

func (p *PoolWrapper) FeeString() string {
	feePercent := float64(p.Fee) * poolFeeToPercent
	return strconv.FormatFloat(feePercent, 'f', -1, 64)
}

func (p *PoolWrapper) String() string {
	out := fmt.Sprintf("%v(chainId=%v fee=%v%%)", p.Symbol(), p.ChainId, p.FeeString())
	return out
}
