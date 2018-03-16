package http

import "github.com/KyberNetwork/reserve-data/common"

type Huobi interface {
	PendingIntermediateTxs(timepoint uint64) (map[common.ActivityID]common.TXEntry, error)
}
