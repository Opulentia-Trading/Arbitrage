package uniswapV2Handler

import "github.com/Opulentia-Trading/Arbitrage/platformListener"

type pairWrapper struct {
	pairAddress string
	token0      *platformListener.ERC20Token
	token1      *platformListener.ERC20Token
}

// ========== Pair Definitions ==========
var (
	// USDC/ETH
	pair_USDC_ETH = pairWrapper{
		pairAddress: "0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc",
		token0:      &platformListener.ERC20TokenUSDC,
		token1:      &platformListener.ERC20TokenWETH,
	}

	// LINK/ETH
	pair_LINK_ETH = pairWrapper{
		pairAddress: "0xa2107FA5B38d9bbd2C461D6EDf11B11A50F6b974",
		token0:      &platformListener.ERC20TokenLINK,
		token1:      &platformListener.ERC20TokenWETH,
	}

	// Map holding references to all defined pairs
	uniswapV2PairMap = map[string]*pairWrapper{
		"USDC/ETH": &pair_USDC_ETH,
		"LINK/ETH": &pair_LINK_ETH,
	}
)
