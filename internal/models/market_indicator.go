package models

type MarketIndicator struct {
	ID            int64
	Pair          string
	TimestampUnix int64
	MMS20         *float64
	MMS50         *float64
	MMS200        *float64
}
