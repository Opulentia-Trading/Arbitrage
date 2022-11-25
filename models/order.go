package models

import (
	"math/big"
	"time"
)

type Order struct {
	Exchange         *Exchange
	Base             string
	Quote            string
	Action           Action
	Price            *big.Rat
	Quantity         *big.Int
	LiqPoolAmountIn  *big.Int // in wei
	LiqPoolAmountOut *big.Int // in wei
	Deadline         time.Duration
	Next             *Order
}
