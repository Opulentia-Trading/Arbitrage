package platformListener

type TickerPrice struct {
	// TODO: Add timestamps
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// Interface defining required methods for a DEX/CEX handler
type ListenerHandler interface {
	TestConnection()
	FetchTickerPriceAll() []TickerPrice
}

type Listener struct {
	PlatformName string
	// supportedTickers []string
	ListenerHandler
}

func NewListener(platformName string, handler ListenerHandler) *Listener {
	return &Listener{
		PlatformName:    platformName,
		ListenerHandler: handler,
	}
}
