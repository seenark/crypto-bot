package bot

import (
	"fmt"
	"mybot/binance"
	"mybot/binance/timeframe"
	"mybot/bot/signal"
	"mybot/candle"
	"mybot/constants"
	"mybot/logs"
	"mybot/mixins"
	"time"

	"github.com/markcheno/go-talib"
)

var profit float64 = 0
var positionQnt float64 = 0
var entryPrice float64 = 0
var assets float64 = 100
var multiply float64 = 50
var useAssets float64 = assets / 10 * multiply
var entrySide string = ""
var lg *logs.Log

// var targetPercent float64 = 0.6
const EMA_CROSS = "EMACross"

type EMACross struct {
	StrategyName string
}

func NewEMACross() *EMACross {
	return &EMACross{
		StrategyName: EMA_CROSS,
	}
}

func (e *EMACross) EMACrossSignal(pair string, klines []binance.KLine) signal.SignalStruct {

	candles := binance.KLinesToCandles(klines)
	cnds := *candles
	// fmt.Printf("cnds: %v\n", cnds)
	haCandles := candle.CandlesToHeikinAshies(*candles)

	slowEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 6)
	// mixins.PrintPretty(slowEma)
	fastEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 1)
	// mixins.PrintPretty(fastEma)

	crossup := talib.Crossover(fastEma, slowEma)
	crossdown := talib.Crossunder(fastEma, slowEma)
	fmt.Printf("crossup: %v\n", crossup)
	lastestCandle := cnds[len(cnds)-1]
	sn := signal.NewSignal(crossup, crossdown, crossdown, crossup, &lastestCandle, EMA_CROSS, pair)
	return *sn
}

func EMACrossProcess(b *binance.BinanceClient, pair string) {
	lg = logs.NewLog("trades.log")
	lg.WriteText(fmt.Sprintf("Start: %s", time.Now().UTC().Format(time.RFC3339)))

	var shift time.Duration = 20
	shift2 := 20
	for {
		end := time.Now()
		fmt.Printf("end: %v\n", end)
		// start := end.AddDate(0, 0, -1)
		start := end.Add((-1 * shift) * 1 * time.Minute)
		// fmt.Printf("start: %v\n", start)

		sec := end.Second()
		fmt.Printf("sec: %v\n", sec)
		if sec == 00 {
			klineCandle := b.GetKLine(pair, timeframe.M1, mixins.ToMilliseconds(start), mixins.ToMilliseconds(end), shift2)
			candles := binance.KLinesToCandles(klineCandle)
			cnds := *candles
			// fmt.Printf("cnds: %v\n", cnds)
			haCandles := candle.CandlesToHeikinAshies(*candles)

			slowEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 6)
			// mixins.PrintPretty(slowEma)
			fastEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 1)
			// mixins.PrintPretty(fastEma)

			long := false
			short := false
			crossup := talib.Crossover(fastEma, slowEma)
			fmt.Printf("crossup: %v\n", crossup)
			lastestCandle := cnds[len(cnds)-1]

			if crossup {
				long = true
				fmt.Printf("long: %v\n", long)

				mixins.PrintPretty(lastestCandle)
				// exitAt0_6Percent(&lastestCandle)
				exit(&lastestCandle, "Exit LONG")
				entry(&lastestCandle, constants.GetShortText(), pair)
			}

			crossdown := talib.Crossunder(fastEma, slowEma)
			fmt.Printf("crossdown: %v\n", crossdown)
			if crossdown {
				short = true
				fmt.Printf("short: %v\n", short)

				mixins.PrintPretty(lastestCandle)
				// exitAt0_6Percent(&lastestCandle)
				exit(&lastestCandle, "Exit SHORT")
				entry(&lastestCandle, constants.GetShortText(), pair)

			}
			time.Sleep(55 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

// func exitAt0_6Percent(lastestCandle *candle.Candle) {
// 	targetPrice := 0.0
// 	if entrySide == LONG {
// 		targetPrice = entryPrice * (1 + (targetPercent / 100))
// 		if lastestCandle.Close >= targetPrice {
// 			exit(lastestCandle, "HIT 0.6% on top")
// 		}
// 	} else if entrySide == SHORT {
// 		targetPrice = entryPrice * (1 - (targetPercent / 100))
// 		if lastestCandle.Close <= targetPrice {
// 			exit(lastestCandle, "HIT 0.6% on bottom")
// 		}
// 	}
// }

func exit(lastestCandle *candle.Candle, comment string) float64 {
	currentProfit := 0.0
	if positionQnt != 0 && entryPrice != 0 {
		// exit
		before := positionQnt * entryPrice
		after := positionQnt * lastestCandle.Close

		if entrySide == constants.GetLongText() {
			currentProfit = after - before
			profit += currentProfit
		} else if entrySide == constants.GetShortText() {
			currentProfit = before - after
			profit += currentProfit
		}
		msg := fmt.Sprintf("EXIT: Profit: %f \t Total Profit: %f \t -- Comment: %s --", currentProfit, profit, comment)
		// go line.SendMessage(msg)
		go lg.WriteText(msg)
		// reset
		entryPrice = 0
		positionQnt = 0
		entrySide = ""
	}
	return currentProfit
}

func entry(lastestCandle *candle.Candle, side string, pair string) {
	entryPrice = lastestCandle.Close
	positionQnt = useAssets / lastestCandle.Close
	entrySide = side
	newTime := time.Unix(int64(lastestCandle.CloseTime)/1000, 0)
	// fmt.Printf("Enter %s, At Price: %f, Qty: %f \n", side, entryPrice, positionQnt)

	msg := fmt.Sprintf("%s \t 🔻🔻🔻🔻🔻 \t ENTER %s \t 🔻🔻🔻🔻🔻 \t - Price: %f \t - Quantity: %f \t - Time: %s", pair, side, lastestCandle.Close, positionQnt, newTime.UTC().Format(time.RFC3339))
	// go line.SendMessage(msg)
	go lg.WriteText(msg)
}
