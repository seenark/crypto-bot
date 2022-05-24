package helpers

import (
	"fmt"
	"runtime"
	"time"
)

const (
	M1  = 1 * time.Minute
	M5  = 5 * time.Minute
	M30 = 30 * time.Minute
	H1  = 1 * time.Hour
	H4  = 4 * time.Hour
	D1  = 24 * time.Hour
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func TimeLast24H() (start, end int) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -1)
	start = ToMilliseconds(startDate)
	end = ToMilliseconds(now)
	return
}

func GetStartEndDate(td time.Duration) (start, end int) {
	now := time.Now()
	startDate := now.Add(-td)
	start = ToMilliseconds(startDate)
	end = ToMilliseconds(now)
	return
}

func ToMilliseconds(t time.Time) int {
	return int(t.UnixNano()) / 1e6
}
