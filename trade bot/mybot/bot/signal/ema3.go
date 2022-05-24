package signal

import (
	"fmt"
	"mybot/binance"
	"mybot/binance/timeframe"
	"mybot/candle"
	"mybot/constants"
	"mybot/mixins"
	"mybot/repository/users"
	"time"

	"github.com/markcheno/go-talib"
)

const (
	EMA3_NAME = "EMA3"
)

type EMA3 struct {
	StrategyName string `json:"strategyName"`
}

func (e *EMA3) Backtest(fakeUser *users.User, largeKlines []binance.KLine) {
	minimumCandles := 100
	for i := 0; i < len(largeKlines)-minimumCandles+1; i++ {
		fmt.Printf("i: %v\n", i)
		var firstCandles [100]binance.KLine
		copy(firstCandles[:], largeKlines[i:i+minimumCandles-1])
		sn := GetEMA3CrossSignal("1000SHIBUSDT", firstCandles)
		if sn.Long {
			fmt.Printf("%s: Go Long\n", sn.Pair)
			fakeUser.EnterPosition(sn.Pair, *sn.LatestCandle, constants.LONG)
		} else {
			fmt.Printf("%s: No Long\n", sn.Pair)
		}
		if sn.CloseLong {
			fmt.Printf("%s: Close Long\n", sn.Pair)
			fakeUser.ExitPosition(sn.Pair, *sn.LatestCandle, constants.LONG)
		} else {
			fmt.Printf("%s: No Close Long\n", sn.Pair)
		}
		if sn.Short {
			fmt.Printf("%s: Go Short\n", sn.Pair)
			fakeUser.EnterPosition(sn.Pair, *sn.LatestCandle, constants.SHORT)
		} else {
			fmt.Printf("%s: No Short\n", sn.Pair)
		}
		if sn.CloseShort {
			fmt.Printf("%s: Close Short\n", sn.Pair)
			fakeUser.ExitPosition(sn.Pair, *sn.LatestCandle, constants.SHORT)
		} else {
			fmt.Printf("%s: No Close Short\n", sn.Pair)
		}
	}
	mixins.PrintPretty(fakeUser)
}

func EMA3CrossProcess(b *binance.BinanceClient, pair string, signalCh chan SignalStruct) {
	var shift time.Duration = 100
	shift2 := 100

	for {
		end := time.Now()
		// fmt.Printf("end: %v\n", end)
		// start := end.AddDate(0, 0, -1)
		start := end.Add((-1 * shift) * 1 * time.Minute)
		// fmt.Printf("start: %v\n", start)

		sec := end.Second()
		fmt.Printf("Pair: %s, sec: %v\n", pair, sec)
		// || sec == 10 || sec == 20 || sec == 30 || sec == 40 || sec == 50
		if sec == 00 {
			klineCandle := b.GetKLine(pair, timeframe.M1, mixins.ToMilliseconds(start), mixins.ToMilliseconds(end), shift2)
			var candles100 [100]binance.KLine
			copy(candles100[:], klineCandle[:])
			sn := GetEMA3CrossSignal(pair, candles100)
			// candles := binance.KLinesToCandles(klineCandle)
			// cnds := *candles
			// haCandles := candle.CandlesToHeikinAshies(*candles)

			// verySlowEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 30)
			// slowEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 6)
			// fastEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 1)
			// lastestCandle := cnds[len(cnds)-1]
			// // exit
			// // exitProcess(pair, &lastestCandle, fastEma, slowEma)

			// currentFastEma := fastEma[len(fastEma)-1]
			// currentSlowEma := slowEma[len(slowEma)-1]
			// currentVerySlowEma := verySlowEma[len(verySlowEma)-1]
			// // entry
			// // enterProcess(pair, &lastestCandle, currentFastEma, currentSlowEma, currentVerySlowEma)
			// sn.Long = canLong(currentFastEma, currentSlowEma, currentVerySlowEma)
			// sn.CloseLong = exitLong(fastEma, slowEma)
			// sn.Short = canShort(currentFastEma, currentSlowEma, currentVerySlowEma)
			// sn.CloseShort = exitShort(fastEma, slowEma)
			// sn.LatestCandle = &lastestCandle
			// sn.Pair = pair
			// sn.StrategyName = EMA3_NAME
			signalCh <- *sn

			time.Sleep(55 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

func GetEMA3CrossSignal(pair string, klineCandles [100]binance.KLine) *SignalStruct {
	candles := binance.KLinesToCandles(klineCandles[0:99])
	cnds := *candles
	haCandles := candle.CandlesToHeikinAshies(*candles)

	verySlowEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 30)
	slowEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 6)
	fastEma := talib.Ema(candle.GetHACloseSeries(*haCandles), 1)
	lastestCandle := cnds[len(cnds)-1]
	// exit
	// exitProcess(pair, &lastestCandle, fastEma, slowEma)

	currentFastEma := fastEma[len(fastEma)-1]
	currentSlowEma := slowEma[len(slowEma)-1]
	currentVerySlowEma := verySlowEma[len(verySlowEma)-1]
	// entry
	long := canLong(currentFastEma, currentSlowEma, currentVerySlowEma)
	closeLong := exitLong(fastEma, slowEma)
	short := canShort(currentFastEma, currentSlowEma, currentVerySlowEma)
	closeShort := exitShort(fastEma, slowEma)
	sn := NewSignal(long, closeLong, short, closeShort, &lastestCandle, pair, EMA3_NAME)

	return sn
}

func canLong(currentFastEma float64, currentSlowEma float64, currentVerySlowEma float64) bool {
	long := false
	if currentFastEma > currentSlowEma && currentFastEma > currentVerySlowEma {
		long = true
	}
	return long
}

func exitLong(fastEma []float64, slowEma []float64) bool {
	return talib.Crossunder(fastEma, slowEma)
}

func canShort(currentFastEma float64, currentSlowEma float64, currentVerySlowEma float64) bool {
	short := false
	if currentFastEma < currentSlowEma && currentFastEma < currentVerySlowEma {
		short = true
	}
	return short
}

func exitShort(fastEma []float64, slowEma []float64) bool {
	return talib.Crossover(fastEma, slowEma)
}

// func enterProcess(pair string, latestCandle *candle.Candle, curFastEma float64, curSlowEma float64, curVerySlowEma float64) {
// 	long := canLong(curFastEma, curSlowEma, curVerySlowEma)

// 	if long {
// 		long = true
// 		fmt.Printf("long: %v\n", long)

// 		mixins.PrintPretty(latestCandle)
// 	} else {
// 		fmt.Println("No Long")
// 	}

// 	short := canShort(curFastEma, curSlowEma, curVerySlowEma)
// 	if short {
// 		fmt.Printf("short: %v\n", short)

// 		mixins.PrintPretty(latestCandle)

// 	} else {
// 		fmt.Println("No Short")
// 	}
// }

// func exitProcess(pair string, lastestCandle *candle.Candle, fastEma []float64, slowEma []float64) {

// 	if exitLong(fastEma, slowEma) {
// 		lg.WriteText("Going to Exit Long")
// 		exit(lastestCandle, "Exit Long", pair)
// 	} else {
// 		fmt.Println("No Exit Long")
// 		lg.WriteText("No Exit Long")
// 	}

// 	if exitShort(fastEma, slowEma) {
// 		lg.WriteText("Going to Exit Short")
// 		exit(lastestCandle, "Exit Short", pair)
// 	} else {
// 		fmt.Println("No Exit Short")
// 		lg.WriteText("No Exit Short")
// 	}
// }

// func exit(lastestCandle *candle.Candle, comment string, pair string) float64 {
// 	currentProfit := 0.0
// 	if positionQnt != 0 && entryPrice != 0 {
// 		// exit
// 		before := positionQnt * entryPrice
// 		after := positionQnt * lastestCandle.Close

// 		if entrySide == LONG {
// 			currentProfit = after - before
// 			profit += currentProfit
// 		} else if entrySide == SHORT {
// 			currentProfit = before - after
// 			profit += currentProfit
// 		}
// 		msg := fmt.Sprintf("EXIT:%s Profit: %f \t Total Profit: %f \t -- Comment: %s --", pair, currentProfit, profit, comment)
// 		go line.SendMessage(msg)
// 		go lg.WriteText(msg)
// 		// reset
// 		entryPrice = 0
// 		positionQnt = 0
// 		entrySide = ""
// 		isInPosition = false
// 	}
// 	return currentProfit
// }

// func entry(lastestCandle *candle.Candle, side string, pair string) {
// 	entryPrice = lastestCandle.Close
// 	positionQnt = useAssets / lastestCandle.Close
// 	entrySide = side
// 	newTime := time.Unix(int64(lastestCandle.CloseTime)/1000, 0)
// 	// fmt.Printf("Enter %s, At Price: %f, Qty: %f \n", side, entryPrice, positionQnt)
// 	isInPosition = true
// 	msg := fmt.Sprintf("%s \t 🔻🔻🔻🔻🔻 \t ENTER %s \t 🔻🔻🔻🔻🔻 \t - Price: %f \t - Quantity: %f \t - Time: %s", pair, side, lastestCandle.Close, positionQnt, newTime.UTC().Format(time.RFC3339))
// 	go line.SendMessage(msg)
// 	go lg.WriteText(msg)
// }

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
