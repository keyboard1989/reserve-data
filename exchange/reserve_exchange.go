package exchange

import "github.com/KyberNetwork/reserve-data/common"

type ReserveExchange struct {
	exchangeStorage ExchangeStorage
}

func NewExchange(storage ExchangeStorage) *ReserveExchange {
	return &ReserveExchange{
		exchangeStorage: storage,
	}
}

func (self *ReserveExchange) GetTradeHistory(fromTime, toTime uint64) (common.AllTradeHistory, error) {
	data, err := self.exchangeStorage.GetTradeHistory(fromTime, toTime)
	return data, err
}
