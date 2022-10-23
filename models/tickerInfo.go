package models

import "time"

type TickerInfo struct {
	Symbol         string `json:"symbol"`
	Base           string
	Quote          string
	Price          string `json:"price"`
	MakerComission string
	TakerComission string
	Timestamp      time.Time
}
