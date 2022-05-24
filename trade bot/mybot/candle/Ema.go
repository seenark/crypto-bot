package candle

import "errors"

type Ema struct {
	Value float64
	Time  float64
}

func MakeEmaSeries(ha []HeikinAshiCandle, emas []float64) ([]Ema, error) {
	emaSeries := []Ema{}
	if len(ha) != len(emas) {
		return nil, errors.New("both slice is not same size")
	}
	for i, h := range ha {
		ema := Ema{
			Value: emas[i],
			Time:  h.CloseTime,
		}
		emaSeries = append(emaSeries, ema)
	}
	return emaSeries, nil
}
