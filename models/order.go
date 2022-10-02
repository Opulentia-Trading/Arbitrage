package models

import (
	"math/big"
	"time"
)

type Order struct {
	Exchange  *Exchange
	Base      string
	Quote     string
	Action    *Action
	Price     *big.Rat
	AmountIn  *big.Int
	AmountOut *big.Int
	Deadline  time.Time
	Next      *Order
}
