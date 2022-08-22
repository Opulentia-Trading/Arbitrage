package platformListener

// Info for an ERC20 token (Ethereum)
type ERC20Token struct {
	Address  string
	Decimals int64
	Symbol   string
	Name     string
}

// ========== Token Definitions ==========
var (
	ERC20TokenUSDC = ERC20Token{
		Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Decimals: 6,
		Symbol:   "USDC",
		Name:     "USD Coin",
	}

	ERC20TokenWETH = ERC20Token{
		Address:  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Decimals: 18,
		Symbol:   "ETH",
		Name:     "Wrapped Ether",
	}

	ERC20TokenWBTC = ERC20Token{
		Address:  "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
		Decimals: 8,
		Symbol:   "WBTC",
		Name:     "Wrapped BTC",
	}

	ERC20TokenAPE = ERC20Token{
		Address:  "0x4d224452801ACEd8B2F0aebE155379bb5D594381",
		Decimals: 18,
		Symbol:   "APE",
		Name:     "ApeCoin",
	}

	ERC20TokenLINK = ERC20Token{
		Address:  "0x514910771AF9Ca656af840dff83E8264EcF986CA",
		Decimals: 18,
		Symbol:   "LINK",
		Name:     "ChainLink Token",
	}
)
