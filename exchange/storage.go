package exchange

import "github.com/KyberNetwork/reserve-data/common"

type ExchangeStorage interface {
	StoreTradeHistory(data common.AllTradeHistory, timepoint uint64) error

	GetTradeHistory(fromTime, toTime uint64) (common.AllTradeHistory, error)
	GetLastIDTradeHistory(exchange, pair string) (string, error)
}
