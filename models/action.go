package models

type Action uint

const (
	BuyShortSpot Action = iota
	BuyLongSpot
	BuyShortFutures
	BuyLongFutures
	SellShortSpot
	SellLongSpot
	SellShortFutures
	SellLongFutures
)

func (a Action) String() string {
	return [...]string{
		"BuyShortSpot",
		"BuyLongSpot",
		"BuyShortFutures",
		"BuyLongFutures",
		"SellShortSpot",
		"SellLongSpot",
		"SellShortFutures",
		"SellLongFutures"}[a]
}
