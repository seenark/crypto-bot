package backtest

import (
	"mybot/bot/signal"
)

type BacktestStrategy interface {
	GetSignal() signal.SignalStruct
}

func Process() {

}
