package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"mybot/binance/timeframe"
	"net/http"
)

const base_url = "https://fapi.binance.com"

type BinanceClient struct {
	BaseUrl    string
	ApiKey     string
	HttpClient *http.Client
}

func NewBinanceClient(apiKey string) *BinanceClient {
	return &BinanceClient{
		BaseUrl:    base_url,
		ApiKey:     apiKey,
		HttpClient: &http.Client{},
	}
}

func (b *BinanceClient) GetMultipleKline(symbol string, timeFrame string, startTimeMilliSec int, endTimeMilliSec int, limit int) []KLine {

	start := startTimeMilliSec
	end := endTimeMilliSec

	nReq := getNumberOfRequest(timeFrame, startTimeMilliSec, endTimeMilliSec)

	return nil
}

func getNumberOfRequest(timeFrame string, startTimeMilliSec int, endTimeMilliSec int) int {

	diff := endTimeMilliSec - startTimeMilliSec
	timeframeMilli := timeframe.GetNumberOfMilliSecondFromTimeFrame(timeFrame)

	nReq := float64(diff) / float64(timeframeMilli)
	return int(math.Ceil(nReq))

}

func (b *BinanceClient) GetKLine(symbol string, interval string, startTime int, endTime int, limit int) []KLine {
	uri := "/fapi/v1/klines"
	url := fmt.Sprintf("%s%s", base_url, uri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	queryString := req.URL.Query()
	queryString.Add("symbol", symbol)
	queryString.Add("interval", interval)
	queryString.Add("startTime", fmt.Sprintf("%d", startTime))
	queryString.Add("endTime", fmt.Sprintf("%d", endTime))
	queryString.Add("limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = queryString.Encode()

	fmt.Println(req.URL)

	if err != nil {
		fmt.Println(err)
	}
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var bd [][]interface{}
	err = json.Unmarshal(body, &bd)
	if err != nil {
		fmt.Println(err)
	}

	newKlines := NewKLine(bd)
	// fmt.Println("newKlines", newKlines)
	return newKlines

}
