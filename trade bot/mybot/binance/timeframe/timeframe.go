package timeframe

import "time"

const (
	M1     = "1m"
	M3     = "3m"
	M5     = "5m"
	M15    = "15m"
	M30    = "30m"
	H1     = "1h"
	H2     = "2h"
	H4     = "4h"
	H6     = "6h"
	H8     = "8h"
	H12    = "12h"
	D1     = "1d"
	D3     = "3d"
	W1     = "1w"
	Month1 = "1M"
)

func GetNumberOfMilliSecondFromTimeFrame(timeframe string) int64 {
	switch timeframe {
	case M1:
		return int64(1 * time.Minute / time.Millisecond)
	case M3:
		return int64(3 * time.Minute / time.Millisecond)
	case M5:
		return int64(5 * time.Minute / time.Millisecond)
	case M15:
		return int64(15 * time.Minute / time.Millisecond)
	case M30:
		return int64(30 * time.Minute / time.Millisecond)
	case H1:
		return int64(1 * time.Hour / time.Millisecond)
	case H2:
		return int64(2 * time.Hour / time.Millisecond)
	case H4:
		return int64(4 * time.Hour / time.Millisecond)
	case H6:
		return int64(6 * time.Hour / time.Millisecond)
	case H8:
		return int64(8 * time.Hour / time.Millisecond)
	case H12:
		return int64(12 * time.Hour / time.Millisecond)
	case D1:
		return int64(1 * time.Hour * 24 / time.Millisecond)
	case D3:
		return int64(3 * time.Hour * 24 / time.Millisecond)
	case W1:
		return int64(1 * time.Hour * 24 * 7 / time.Millisecond)
	case Month1:
		return int64(1 * time.Hour * 24 * 7 * 30 / time.Millisecond) // this may be incorrect in month that have 31 days
	default:
		return 0
	}

}
