package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iamjinlei/go-tart"
	"github.com/seenark/rebalance-bot/binance"
	"github.com/seenark/rebalance-bot/helpers"
)

type Account struct {
	StartDeposit float64
	CoinSymbol   string
	CoinRatio    float64
	CoinAmount   float64
	CoinValue    float64
	UsdtAmount   float64
	topAtrPrice  float64
	lowAtrPrice  float64
}

func (acc Account) UsdtRatio() (ratio float64) {
	ratio = 1 - acc.CoinRatio
	return
}

var account = Account{
	StartDeposit: 1000,
	CoinSymbol:   "BTCUSDT",
	CoinRatio:    0.5,
	CoinAmount:   0,
	CoinValue:    0,
	UsdtAmount:   0,
	topAtrPrice:  0,
	lowAtrPrice:  0,
}

func main() {
	// rebalance(&account, 100)
	// rebalance(&account, 110)

	binanceClient := binance.NewBinanceClient(&http.Client{})
	chartBarNo := 500
	start, end := helpers.GetStartEndDate(time.Duration(chartBarNo) * helpers.M30)
	kline := binanceClient.GetKLine(account.CoinSymbol, "30m", start, end, chartBarNo)
	// fmt.Printf("kline: %v\n", kline)
	atr := tart.NewAtr(14)
	for index, kl := range kline {
		atrValue := atr.Update(kl.High, kl.Low, kl.Close)
		// fmt.Printf("atrValue: %v\n", atrValue)
		// fmt.Printf("#%d -> price: %f -> ", index, kl.Close)
		if (index > 14) && (kl.Close <= account.lowAtrPrice || kl.Close >= account.topAtrPrice) {
			rebalance(&account, kl.Close, atrValue)
		}
	}

}

func rebalance(acc *Account, currentPrice float64, atrValue float64) {
	amountToBuy := 0.0
	amountToSell := 0.0

	if acc.StartDeposit <= 0 {
		// fmt.Printf("acc: %v\n", acc)
		currentAmount := acc.CoinAmount
		currentValue := currentAmount * currentPrice
		sum := currentValue + acc.UsdtAmount

		acc.CoinValue = sum * acc.CoinRatio
		acc.UsdtAmount = sum * acc.UsdtRatio()

		newCoinAmount := (acc.CoinValue / currentPrice)
		// fmt.Printf("newCoinAmount: %v\n", newCoinAmount)
		if newCoinAmount < acc.CoinAmount {
			amountToSell = acc.CoinAmount - newCoinAmount
		} else {
			amountToBuy = newCoinAmount - acc.CoinAmount
		}

		_ = amountToBuy
		_ = amountToSell
		acc.topAtrPrice = currentPrice + atrValue
		acc.lowAtrPrice = currentPrice - atrValue
		// fmt.Printf("amountToBuy: %v\n", amountToBuy)
		// fmt.Printf("amountToSell: %v\n", amountToSell)
		newSum := acc.CoinValue + acc.UsdtAmount
		fmt.Printf("newSum: %v\n", newSum)
	} else {
		acc.CoinValue = acc.CoinRatio * acc.StartDeposit
		newUsdtValue := acc.UsdtRatio() * acc.StartDeposit
		acc.CoinAmount = acc.CoinValue / currentPrice
		acc.UsdtAmount = newUsdtValue
		acc.StartDeposit -= (acc.CoinValue + newUsdtValue)
		acc.topAtrPrice = currentPrice + atrValue
		acc.lowAtrPrice = currentPrice - atrValue
		// fmt.Printf("acc: %v\n", acc)
	}
}
