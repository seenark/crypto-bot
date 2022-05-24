package binance

import (
	"fmt"
	"mybot/candle"
	"strconv"
)

type KLine struct {
	OpenTime                float64
	Open                    float64
	High                    float64
	Low                     float64
	Close                   float64
	Volume                  float64
	CloseTime               float64
	QuoteAssetVolume        float64
	NumberOfTrades          float64
	TakerBuyAssetVolume     float64
	TakerBuyQuoteAssetVolum float64
	Ignore                  float64

	// 1499040000000,      // Open time
	//   "0.01634790",       // Open
	//   "0.80000000",       // High
	//   "0.01575800",       // Low
	//   "0.01577100",       // Close
	//   "148976.11427815",  // Volume
	//   1499644799999,      // Close time
	//   "2434.19055334",    // Quote asset volume
	//   308,                // Number of trades
	//   "1756.87402397",    // Taker buy base asset volume
	//   "28.46694368",      // Taker buy quote asset volume
	//   "17928899.62484339" // Ignore.
}

func NewKLine(klineSlice [][]interface{}) []KLine {
	// fmt.Println("v", klineSlice)
	klines := []KLine{}

	for _, kl := range klineSlice {
		kline := KLine{}
		for i, k := range kl {
			str := fmt.Sprintf("%v", k)
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				appendKline(i, f, &kline)
			}
		}

		// fmt.Println("kline", kline)
		klines = append(klines, kline)
	}
	return klines
}

func appendKline(indexNumber int, f float64, kl *KLine) *KLine {
	switch indexNumber {
	case 0:
		kl.OpenTime = f
	case 1:
		kl.Open = f
	case 2:
		kl.High = f
	case 3:
		kl.Low = f
	case 4:
		kl.Close = f
	case 5:
		kl.Volume = f
	case 6:
		kl.CloseTime = f
	case 7:
		kl.QuoteAssetVolume = f
	case 8:
		kl.NumberOfTrades = f
	case 9:
		kl.TakerBuyAssetVolume = f
	case 10:
		kl.TakerBuyQuoteAssetVolum = f
	case 11:
		kl.Ignore = f

	}

	return kl
}

func (kl *KLine) MakeCandleFromKLine() *candle.Candle {
	return candle.NewCandle(kl.Open, kl.High, kl.Low, kl.Close, kl.CloseTime)
}

func KLinesToCandles(klines []KLine) *[]candle.Candle {
	candles := []candle.Candle{}
	for _, kl := range klines {
		candles = append(candles, *kl.MakeCandleFromKLine())
	}
	return &candles
}
