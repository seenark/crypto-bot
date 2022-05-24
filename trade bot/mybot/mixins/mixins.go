package mixins

import (
	"encoding/json"
	"fmt"
	"time"
)

func Max(numbers ...float64) float64 {
	var max float64 = 0
	for _, num := range numbers {
		if num > max {
			max = num
		}
	}
	return max
}

func Min(numbers ...float64) float64 {
	var min float64 = 999999999
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	return min
}

func ToMilliseconds(t time.Time) int {
	return int(t.UnixNano()) / 1e6
}

func PrintType(x interface{}) {
	xType := fmt.Sprintf("%T", x)
	fmt.Println(xType) // "[]int"
}

func PrintPretty(x interface{}) {
	j, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("Pretty Print Error", err)
	}
	fmt.Println(string(j))
}

func PrintTime(milliSecond float64) {
	newTime := time.Unix(int64(milliSecond)/1000, 0)
	fmt.Printf(" Time: %v\n", newTime.UTC().Format(time.RFC3339))
}
