package metric

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type MetricStorage interface {
	StoreMetric(data *MetricEntry, timepoint uint64) error
	StoreTokenTargetQty(id, data string) error
	StorePendingTargetQty(data, dataType string) error

	GetMetric(tokens []common.Token, fromTime, toTime uint64) (map[string]MetricList, error)
	GetTokenTargetQty() (TokenTargetQty, error)
	GetPendingTargetQty() (TokenTargetQty, error)

	RemovePendingTargetQty() error
}
