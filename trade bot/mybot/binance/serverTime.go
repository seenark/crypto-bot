package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetServerTime(b *BinanceClient) time.Time {
	uri := "/fapi/v1/time"
	url := fmt.Sprintf("%s%s", base_url, uri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
	var serverTime map[string]int64
	err = json.Unmarshal(body, &serverTime)
	if err != nil {
		fmt.Println(err)
	}

	unix := serverTime["serverTime"] / 1000
	t := time.Unix(unix, 0)

	fmt.Printf("serverTime: %v\n", t)
	return t
}
