package uniswapV3Handler

import "github.com/Opulentia-Trading/Arbitrage/platformListener"

type poolWrapper struct {
	poolAddress string
	token0      *platformListener.ERC20Token
	token1      *platformListener.ERC20Token
}

// ========== Pool Definitions ==========
var (
	// USDC/ETH
	pool_USDC_ETH = poolWrapper{
		poolAddress: "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
		token0:      &platformListener.ERC20TokenUSDC,
		token1:      &platformListener.ERC20TokenWETH,
	}

	// WBTC/ETH
	pool_WBTC_ETH = poolWrapper{
		poolAddress: "0xCBCdF9626bC03E24f779434178A73a0B4bad62eD",
		token0:      &platformListener.ERC20TokenWBTC,
		token1:      &platformListener.ERC20TokenWETH,
	}

	// APE/ETH
	pool_APE_ETH = poolWrapper{
		poolAddress: "0xAc4b3DacB91461209Ae9d41EC517c2B9Cb1B7DAF",
		token0:      &platformListener.ERC20TokenAPE,
		token1:      &platformListener.ERC20TokenWETH,
	}

	// LINK/ETH
	pool_LINK_ETH = poolWrapper{
		poolAddress: "0xa6Cc3C2531FdaA6Ae1A3CA84c2855806728693e8",
		token0:      &platformListener.ERC20TokenLINK,
		token1:      &platformListener.ERC20TokenWETH,
	}

	// Map holding references to all defined pools
	uniswapV3PoolMap = map[string]*poolWrapper{
		"USDC/ETH": &pool_USDC_ETH,
		"WBTC/ETH": &pool_WBTC_ETH,
		"APE/ETH":  &pool_APE_ETH,
		"LINK/ETH": &pool_LINK_ETH,
	}
)
