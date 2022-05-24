package signal

import "mybot/candle"

const (
	LONG  = "LONG"
	SHORT = "SHORT"
)

type SignalStruct struct {
	Long         bool           `json:"long"`
	CloseLong    bool           `json:"closeLong"`
	Short        bool           `json:"short"`
	CloseShort   bool           `json:"closeShort"`
	LatestCandle *candle.Candle `json:"latestCandle"`
	Pair         string         `json:"pair"`
	StrategyName string         `json:"strategyName"`
}

func NewSignal(long bool, closeLong bool, short bool, closeShort bool, latestCandle *candle.Candle, pair string, strategyName string) *SignalStruct {
	return &SignalStruct{
		Long:         long,
		CloseLong:    closeLong,
		Short:        short,
		CloseShort:   closeShort,
		LatestCandle: latestCandle,
		Pair:         pair,
		StrategyName: strategyName,
	}
}
