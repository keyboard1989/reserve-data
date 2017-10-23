package reserve

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type TestReserveData struct {
}

func (self TestReserveData) CurrentPriceVersion() common.Version {
	return common.Version(0)
}

func (self TestReserveData) GetAllPrices() (common.AllPriceResponse, error) {
	return common.AllPriceResponse{}, nil
}

func (self TestReserveData) GetOnePrice(common.TokenPairID) (common.OnePriceResponse, error) {
	return common.OnePriceResponse{}, nil
}

func (self TestReserveData) Run() error {
	return nil
}

func NewTestReserveData() *TestReserveData {
	return &TestReserveData{}
}
