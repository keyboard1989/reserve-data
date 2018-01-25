package metric

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type MetricStorage interface {
	StoreMetric(data *MetricEntry, timepoint uint64) error
	StoreTokenTargetQty(data TokenTargetQty) error

	GetMetric(tokens []common.Token, fromTime, toTime uint64) (map[string]MetricList, error)
	GetTokenTargetQty() (map[string]TargetQty, error)
}
