package platformListener

type TickerPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type ListenerHandler interface {
	PingTest()
	FetchTickerPriceAll() []TickerPrice
}

type Listener struct {
	platformName     string
	supportedTickers []string
	ListenerHandler
}

func NewListener(platformName string, supportedTickers []string) *Listener {
	return &Listener{platformName: platformName, supportedTickers: supportedTickers}
}

// Bind a DEX/CEX handler to a listener instance
func (l *Listener) BindHandler(handler ListenerHandler) {
	l.ListenerHandler = handler
}
