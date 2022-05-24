package bot

import (
	"fmt"
	"mybot/binance"
	"mybot/bot/signal"
	"mybot/constants"
	"mybot/mixins"
	"mybot/repository/users"
)

type RunningStrategy struct {
	WatchUsers map[string][]*users.User
}

func BotProcess(signalCh chan signal.SignalStruct, runningStrategies map[string]RunningStrategy) {

	for {
		sn := <-signalCh
		fmt.Println("block")
		mixins.PrintPretty(sn)
		go handleSignal(sn, runningStrategies)
		fmt.Println("Un block")
	}
}

func isSignalPairWatchedByUserPairList(pair string, u *users.User) bool {
	isInUserWatchPair := false
	fmt.Printf("u.WatchPairs: %v\n", u.WatchPairs)
	for _, p := range u.WatchPairs {
		if p == pair {
			isInUserWatchPair = true
			break
		}
	}
	return isInUserWatchPair
}

func handleSignal(sn signal.SignalStruct, runningStrategies map[string]RunningStrategy) {
	fmt.Printf("sn: %v\n", sn)
	fmt.Printf("runningStrategies: %v\n", runningStrategies)
	userList := runningStrategies[sn.StrategyName].WatchUsers[sn.Pair]
	for _, u := range userList {
		if !isSignalPairWatchedByUserPairList(sn.Pair, u) {
			continue
		}
		mixins.PrintPretty(u)
		if sn.Long {
			fmt.Printf("%s: Go Long\n", sn.Pair)
			u.EnterPosition(sn.Pair, *sn.LatestCandle, constants.LONG)
		} else {
			fmt.Printf("%s: No Long\n", sn.Pair)
		}
		if sn.CloseLong {
			fmt.Printf("%s: Close Long\n", sn.Pair)
			u.ExitPosition(sn.Pair, *sn.LatestCandle, constants.LONG)
		} else {
			fmt.Printf("%s: No Close Long\n", sn.Pair)
		}
		if sn.Short {
			fmt.Printf("%s: Go Short\n", sn.Pair)
			u.EnterPosition(sn.Pair, *sn.LatestCandle, constants.SHORT)
		} else {
			fmt.Printf("%s: No Short\n", sn.Pair)
		}
		if sn.CloseShort {
			fmt.Printf("%s: Close Short\n", sn.Pair)
			u.ExitPosition(sn.Pair, *sn.LatestCandle, constants.SHORT)
		} else {
			fmt.Printf("%s: No Close Short\n", sn.Pair)
		}
	}
}

func New3EMASignal(pair string, b *binance.BinanceClient, ch chan signal.SignalStruct, user *users.User, runningStrategies map[string]RunningStrategy) {
	// newStrategy := RunningStrategy{
	// 	Pair:      pair,
	// 	StrgyName: signal.EMA3_NAME,
	// }
	user.WatchPairs = append(user.WatchPairs, pair)
	// newStrategy.WatchUsers = append(newStrategy.WatchUsers, *user)
	if strgy, ok := runningStrategies[signal.EMA3_NAME]; ok {
		fmt.Printf("strgy: %v\n", strgy)
		userList := strgy.WatchUsers[pair]
		userList = append(userList, user)
		runningStrategies[signal.EMA3_NAME].WatchUsers[pair] = userList
		go signal.EMA3CrossProcess(b, pair, ch)
	} else {
		// create new
		runningStrategies[signal.EMA3_NAME] = RunningStrategy{
			WatchUsers: make(map[string][]*users.User),
		}

		userList := strgy.WatchUsers[pair]
		userList = append(userList, user)
		runningStrategies[signal.EMA3_NAME].WatchUsers[pair] = userList
		go signal.EMA3CrossProcess(b, pair, ch)

	}

}
