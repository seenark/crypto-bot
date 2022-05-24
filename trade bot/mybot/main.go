package main

import (
	"fmt"
	"mybot/binance"
	"mybot/bot"
	"mybot/bot/signal"
	"mybot/line"
	"mybot/repository/users"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	signalCh chan signal.SignalStruct

	/** map with strategy name */
	RunningStrategies map[string]bot.RunningStrategy
)

func main() {
	initTimeZone()
	// config := configs.GetConfig()
	// mongo
	// ctx := context.TODO()
	// mgClient := connectMongo(config.Mongo.Username, config.Mongo.Password)
	// db := mgClient.Database(config.Mongo.DbName)
	// emaCollection := db.Collection(config.Mongo.EmaCollection)
	// userCollection := db.Collection(config.Mongo.UserCollection)

	// _ = ctx
	// _ = emaCollection
	// _ = userCollection

	// usrRepo := users.NewUserRepository(userCollection, ctx)
	// userService := usersService.NewUserService(usrRepo)
	hadesgod := users.User{
		Id:              "",
		Name:            "HadesGod",
		WatchPairs:      []string{},
		Equity:          100,
		Position:        &users.Position{},
		Profit:          0,
		Leverage:        50,
		UseAssetPercent: 5,
		Line:            &line.Line{Token: "sBV0pqTRkK5NhYtSUnQHLevNWAVJuBsjRsyPfG5nEMy"},
	}
	_ = hadesgod
	// userService.Create(u)

	/** binance */
	binanceClient := binance.NewBinanceClient("")
	_ = binanceClient
	/*
		backtest
	*/
	// end := time.Now()
	// start := end.Add(-1 * 1000 * time.Minute)
	// klines := binanceClient.GetKLine("1000SHIBUSDT", timeframe.M1, mixins.ToMilliseconds(start), mixins.ToMilliseconds(end), 1000)
	// end = start.Add(-1 * time.Minute)
	// start = end.Add(-1 * 1000 * time.Minute)
	// klines2 := binanceClient.GetKLine("1000SHIBUSDT", timeframe.M1, mixins.ToMilliseconds(start), mixins.ToMilliseconds(end), 1000)
	// fmt.Printf("klines: %v\n", len(klines))

	/* bot processes */
	// pair := "1000SHIBUSDT"
	// go bot.EMACrossProcess(binanceClient, pair)
	// RunningStrategies = make(map[string]bot.RunningStrategy)

	// signalCh = make(chan signal.SignalStruct)
	// go bot.BotProcess(signalCh, RunningStrategies)
	// bot.New3EMASignal("1000SHIBUSDT", binanceClient, signalCh, &hadesgod, RunningStrategies)

	e := echo.New()
	port := 8080
	// e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%d", config.EchoApp.Port)))
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%d", port)))
}

// func connectMongo(username string, password string) *mongo.Client {
// 	clientOptions := options.Client().
// 		ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@hdgcluster.xmgsx.mongodb.net/myFirstDatabase?retryWrites=true&w=majority", username, password))
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return client
// }

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}
