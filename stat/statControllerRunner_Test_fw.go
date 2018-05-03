package stat

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-data/stat/statstoragecontroller"
)

type ControllerRunnerTest struct {
	cr statstoragecontroller.ControllerRunner
}

func NewControllerRunnerTest(controllerRunner statstoragecontroller.ControllerRunner) *ControllerRunnerTest {
	return &ControllerRunnerTest{controllerRunner}
}

func (self *ControllerRunnerTest) TestAnalyticStorageControlTicker(nanosec int64) error {
	if err := self.cr.Start(); err != nil {
		return err
	}
	startTime := time.Now()
	t := <-self.cr.GetAnalyticStorageControlTicker()
	timeTook := t.Sub(startTime).Nanoseconds()
	upperRange := nanosec + nanosec/10
	lowerRange := nanosec - nanosec/10
	if timeTook < lowerRange || timeTook > upperRange {
		return fmt.Errorf("expect ticker in between %d and %d nanosec, but it came in %d instead", lowerRange, upperRange, timeTook)
	}
	if err := self.cr.Stop(); err != nil {
		return err
	}
	return nil
}
