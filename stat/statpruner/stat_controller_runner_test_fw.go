package statpruner

import "log"

type ControllerRunnerTest struct {
	cr ControllerRunner
}

func NewControllerRunnerTest(controllerRunner ControllerRunner) *ControllerRunnerTest {
	return &ControllerRunnerTest{controllerRunner}
}

func (self *ControllerRunnerTest) TestAnalyticStorageControlTicker(nanosec int64) error {
	if err := self.cr.Start(); err != nil {
		return err
	}
	t := <-self.cr.GetAnalyticStorageControlTicker()
	log.Printf("ticked: %s", t.String())
	if err := self.cr.Stop(); err != nil {
		return err
	}
	return nil
}
