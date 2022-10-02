package models

type ExchangeType uint

const (
	Centralized ExchangeType = iota
	Decentralized
)

type Exchange struct {
	Type ExchangeType
	Name string
}

func (t ExchangeType) String() string {
	return [...]string{"CEX", "DEX"}[t]
}
