package exchange

import "github.com/KyberNetwork/reserve-data/common"

type ExchangeStorage interface {
	StoreTradeHistory(data common.AllTradeHistory) error

	GetTradeHistory(fromTime, toTime uint64) (common.AllTradeHistory, error)
	GetLastIDTradeHistory(exchange, pair string) (string, error)
}
