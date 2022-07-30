package uniswapV3Handler

type tokenWrapper struct {
	address  string
	decimals int64
	symbol   string
	name     string
}

type poolWrapper struct {
	poolAddress string
	token0      *tokenWrapper
	token1      *tokenWrapper
}

// ========== Token Definitions ==========
var (
	tokenUSDC = tokenWrapper{
		address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		decimals: 6,
		symbol:   "USDC",
		name:     "USD Coin",
	}

	tokenWETH = tokenWrapper{
		address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		decimals: 18,
		symbol:   "ETH",
		name:     "Wrapped Ether",
	}

	tokenWBTC = tokenWrapper{
		address:  "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
		decimals: 8,
		symbol:   "WBTC",
		name:     "Wrapped BTC",
	}

	tokenAPE = tokenWrapper{
		address:  "0x4d224452801ACEd8B2F0aebE155379bb5D594381",
		decimals: 18,
		symbol:   "APE",
		name:     "ApeCoin",
	}

	tokenLINK = tokenWrapper{
		address:  "0x514910771AF9Ca656af840dff83E8264EcF986CA",
		decimals: 18,
		symbol:   "LINK",
		name:     "ChainLink Token",
	}
)

// ========== Pool Definitions ==========
var (
	// USDC/ETH
	pool_USDC_ETH = poolWrapper{
		poolAddress: "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
		token0:      &tokenUSDC,
		token1:      &tokenWETH,
	}

	// WBTC/ETH
	pool_WBTC_ETH = poolWrapper{
		poolAddress: "0xCBCdF9626bC03E24f779434178A73a0B4bad62eD",
		token0:      &tokenWBTC,
		token1:      &tokenWETH,
	}

	// APE/ETH
	pool_APE_ETH = poolWrapper{
		poolAddress: "0xAc4b3DacB91461209Ae9d41EC517c2B9Cb1B7DAF",
		token0:      &tokenAPE,
		token1:      &tokenWETH,
	}

	// LINK/ETH
	pool_LINK_ETH = poolWrapper{
		poolAddress: "0xa6Cc3C2531FdaA6Ae1A3CA84c2855806728693e8",
		token0:      &tokenLINK,
		token1:      &tokenWETH,
	}

	// Map holding references to all defined pools
	uniswapV3PoolMap = map[string]*poolWrapper{
		"USDC/ETH": &pool_USDC_ETH,
		"WBTC/ETH": &pool_WBTC_ETH,
		"APE/ETH":  &pool_APE_ETH,
		"LINK/ETH": &pool_LINK_ETH,
	}
)
