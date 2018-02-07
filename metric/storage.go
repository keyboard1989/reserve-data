package metric

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type MetricStorage interface {
	StoreMetric(data *MetricEntry, timepoint uint64) error
	StoreTokenTargetQty(id, data string) error
	StorePendingTargetQty(data, dataType string) error
	StoreRebalanceControl(status bool) error

	GetMetric(tokens []common.Token, fromTime, toTime uint64) (map[string]MetricList, error)
	GetTokenTargetQty() (TokenTargetQty, error)
	GetPendingTargetQty() (TokenTargetQty, error)
	GetRebalanceControl() (RebalanceControl, error)

	RemovePendingTargetQty() error
}
