package input

import "github.com/shopspring/decimal"

type QueryData struct {
	Offset        int
	Limit         int
	Id            uint
	Code          string
	Category      string
	PriceLessThan decimal.Decimal
}
