package users

import (
	"fmt"

	"mybot/candle"
	"mybot/constants"
	"mybot/line"
	"time"
)

type Position struct {
	IsInPosition bool
	Side         string
	EntryPrice   float64
	Quantity     float64
}

func (p *Position) Reset() {
	p.EntryPrice = 0
	p.Side = ""
	p.EntryPrice = 0
	p.Quantity = 0
	p.IsInPosition = false
}

type User struct {
	Id              string     `json:"id" bson:"_id"`
	Name            string     `json:"name" bson:"name,omitempty"` // omitempty คือถ้าไม่มีค่าก็ไม่ต้องใส่ key เข้าไปใน db
	WatchPairs      []string   `json:"watchPair" bson:"watchPair"`
	Equity          float64    `json:"equity" bson:"equity,omitempty"`
	Position        *Position  `json:"position" bson:"position"`
	Profit          float64    `json:"profit" bson:"profit"`
	Leverage        float64    `json:"Leverage" bson:"Leverage"`
	UseAssetPercent float64    `json:"useAssetPercent" bson:"useAssetPercent"`
	Line            *line.Line `json:"line" bson:"line"`
}

type UserRepository interface {
	Create(User) (string, error)
	GetAll() ([]User, error)
	GetUserById(string) (*User, error)
	UpdateUserById(string, User) (*User, error)
	DeleteUserById(string) (int, error)
}

func (u *User) EnterPosition(pair string, latestCandle candle.Candle, side string) {
	if u.Position.IsInPosition {
		fmt.Printf("Signal: %s come but we are in %s position\n", side, u.Position.Side)
		return
	}
	useAsset := u.Equity * u.Leverage * (u.UseAssetPercent / 100)
	newTime := time.Unix(int64(latestCandle.CloseTime)/1000, 0)
	u.Position.Side = side
	u.Position.EntryPrice = latestCandle.Close
	u.Position.IsInPosition = true
	u.Position.Quantity = useAsset / latestCandle.Close
	msg := fmt.Sprintf("%s \t 🔻🔻🔻🔻🔻 \t ENTER %s \t 🔻🔻🔻🔻🔻 \t - Price: %f \t - Quantity: %f \t USDT used: %f \t - Time: %s", pair, side, latestCandle.Close, u.Position.Quantity, useAsset, newTime.UTC().Format(time.RFC3339))
	go u.Line.SendMessage(msg)
}

func (u *User) ExitPosition(pair string, lastCandle candle.Candle, side string) {
	if !u.Position.IsInPosition {
		fmt.Printf("You have no any Position\n")
		return
	}
	if side != u.Position.Side {
		msg := "Close signal side is not same in Our position side"
		u.Line.SendMessage(msg)
	}
	before := u.Position.Quantity * u.Position.EntryPrice
	after := u.Position.Quantity * lastCandle.Close
	currentProfit := 0.0
	if side == constants.LONG {
		currentProfit = after - before
	} else {
		currentProfit = before - after
	}
	u.Profit += currentProfit
	u.Position.Reset()
	newTime := time.Unix(int64(lastCandle.CloseTime)/1000, 0)
	msg := fmt.Sprintf("EXIT:%s Profit: %f \t Total Profit: %f \t Time: %s \t-- Comment: %s --", pair, currentProfit, u.Profit, newTime.UTC().Format(time.RFC3339), "")
	go u.Line.SendMessage(msg)
}
