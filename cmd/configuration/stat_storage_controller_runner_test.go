package configuration

import (
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-data/stat/statstoragecontroller"
)

const (
	STATCONTROLLERTESTDURATION time.Duration = time.Second * 1
)

func SetupTickerTestForControllerRunner(duration time.Duration) (*statstoragecontroller.ControllerRunnerTest, error) {
	tickerRuner := statstoragecontroller.NewControllerTickerRunner(duration)
	return statstoragecontroller.NewControllerRunnerTest(tickerRuner), nil
}

func doTickerforControllerRunnerTest(f func(tester *statstoragecontroller.ControllerRunnerTest, t *testing.T), t *testing.T) {
	tester, err := SetupTickerTestForControllerRunner(STATCONTROLLERTESTDURATION)
	if err != nil {
		t.Fatalf("Testing Ticker Runner as Controller Runner: init failed(%s)", err)
	}
	f(tester, t)
}

func TestTickerRunnerForControllerRunner(t *testing.T) {
	doTickerforControllerRunnerTest(func(tester *statstoragecontroller.ControllerRunnerTest, t *testing.T) {
		if err := tester.TestAnalyticStorageControlTicker(STATCONTROLLERTESTDURATION.Nanoseconds()); err != nil {
			t.Fatalf("Testing Ticker Runner ")
		}
	}, t)
}
