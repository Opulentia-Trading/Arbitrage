package validator

import (
	"math"
)

type PlatformInfo struct {
	TYPE string
	NAME string
}
type ArbResult struct {
	SYMBOL        string
	PLATFORM      *PlatformInfo
	ACTION        string
	AMOUNT        float64
	TRANSACTIONID string
	NEXT          *ArbResult
}

type Token struct {
	Name    string
	Reserve float64
}
type Platform struct {
	Tokens [2]*Token
	Name   string
}

func amount_out(to_buy float64, p0t0 float64, p0t1 float64) float64 {
	amount_in_with_fee := to_buy * 997
	numerator := amount_in_with_fee * p0t1
	denominator := (p0t0 * 1000) + amount_in_with_fee
	return numerator / denominator
}

func amount_in(to_buy float64, p0t0 float64, p0t1 float64) float64 {
	numerator := p0t0 * to_buy * 1000
	denominator := (p0t1 - to_buy) * 997
	return (numerator / denominator) + 1
}

func calc_max_to_buy(p0t0 float64, p0t1 float64, p1t0 float64, p1t1 float64) float64 {

	values := []float64{p0t0, p0t1, p1t0, p1t1}
	minimum := values[0]
	for _, v := range values {
		if v < minimum {
			minimum = v
		}
	}
	var d float64
	if minimum > 1e24 {
		d = 1e20
	} else if minimum > 1e23 {
		d = 1e19
	} else if minimum > 1e22 {
		d = 1e18
	} else if minimum > 1e21 {
		d = 1e17
	} else if minimum > 1e19 {
		d = 1e16
	} else if minimum > 1e18 {
		d = 1e15
	} else if minimum > 1e17 {
		d = 1e14
	} else if minimum > 1e16 {
		d = 1e13
	} else if minimum > 1e15 {
		d = 1e12
	} else if minimum > 1e14 {
		d = 1e11
	} else {
		d = 1e10
	}

	p0t0 = p0t0 / d
	p0t1 = p0t1 / d
	p1t0 = p1t0 / d
	p1t1 = p1t1 / d

	term1 := float64(p0t0*p0t1 - p1t0*p1t1)
	term2 := float64(2 * p0t1 * p1t1 * (p0t0 + p1t0))
	term3 := float64(p0t1 * p1t1 * (p0t0*p1t1 - p1t0*p0t1))

	discriminant := math.Pow(term2, 2) - (4 * term1 * term3)

	if discriminant < 0 {
		return -1
	}

	x1 := (-term2 + math.Sqrt(discriminant)) / (2 * term1)
	x2 := (-term2 - math.Sqrt(discriminant)) / (2 * term1)

	if x1 > 0 && x1 < p0t1 && x1 < p1t1 {
		return x1 * d
	} else if x2 > 0 && x2 < p0t1 && x2 < p1t1 {
		return x2 * d
	}

	return -1
}

func main(platform1 Platform, platform2 Platform) *ArbResult {
	pool0token0, pool0token1 := platform1.Tokens[0], platform1.Tokens[1]
	pool1token0, pool1token1 := platform2.Tokens[0], platform2.Tokens[1]

	amount_to_buy := calc_max_to_buy(pool0token0.Reserve, pool0token1.Reserve, pool1token0.Reserve, pool1token1.Reserve)
	if amount_to_buy <= 0 {
		return nil
	}

	amount_to_pay_back := amount_in(amount_to_buy, pool0token0.Reserve, pool0token1.Reserve)

	amount_received := amount_out(amount_to_buy, pool1token1.Reserve, pool1token0.Reserve)
	if amount_received < amount_to_pay_back {
		return nil
	}

	return &ArbResult{
		SYMBOL:        "",
		PLATFORM:      nil,
		ACTION:        "BUY",
		AMOUNT:        0.0,
		TRANSACTIONID: "",
		NEXT:          nil,
	}
}
