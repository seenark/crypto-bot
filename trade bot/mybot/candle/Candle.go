package candle

import (
	"math"
	"sort"
)

type Candle struct {
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	CloseTime float64 `json:"close_time"`
}

func NewCandle(open float64, high float64, low float64, close float64, closeTime float64) *Candle {
	return &Candle{
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		CloseTime: closeTime,
	}
}

type HeikinAshiCandle struct {
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	CloseTime float64 `json:"close_time"`
}

func (c Candle) HeikinAshi(previousCandle HeikinAshiCandle) HeikinAshiCandle {
	close := (c.Open + c.Close + c.High + c.Low) / 4
	open := (previousCandle.Open + previousCandle.Close) / 2
	high := math.Max(c.High, math.Max(open, close))
	// low := mixins.Min(c.Low, open, close)
	low := math.Min(c.Low, math.Min(open, close))

	return HeikinAshiCandle{
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		CloseTime: c.CloseTime,
	}
}

func SortByCloseTime(c []Candle) {
	sort.SliceStable(c, func(i, j int) bool {
		return c[i].CloseTime > c[j].CloseTime
	})
}

func CandlesToHeikinAshies(candles []Candle) *[]HeikinAshiCandle {
	ha := []HeikinAshiCandle{}
	for index, candle := range candles {
		if index == 0 {
			newHa := HeikinAshiCandle(candle)
			ha = append(ha, newHa)
		} else {
			// if index == 1 {
			// 	fmt.Println("index: ", index, ha[index-1])
			// }
			ha = append(ha, candle.HeikinAshi(ha[index-1]))
		}
	}
	return &ha
}

func SortHaCandleByCloseTime(c []HeikinAshiCandle) {
	sort.SliceStable(c, func(i, j int) bool {
		return c[i].CloseTime > c[j].CloseTime
	})
}

func GetHACloseSeries(ha []HeikinAshiCandle) []float64 {
	nums := []float64{}
	for _, c := range ha {
		nums = append(nums, c.Close)
	}
	return nums
}
