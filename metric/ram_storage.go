package metric

import (
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
)

const MAX_CAPACITY int = 1000

type RamMetricStorage struct {
	mu             sync.RWMutex
	data           []*MetricEntry
	tokenTargetQty TokenTargetQty
}

func NewRamMetricStorage() *RamMetricStorage {
	return &RamMetricStorage{
		mu:             sync.RWMutex{},
		data:           []*MetricEntry{},
		tokenTargetQty: TokenTargetQty{},
	}
}

func (self *RamMetricStorage) StoreMetric(data *MetricEntry, timepoint uint64) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.data = append(self.data, data)
	first := len(self.data) - MAX_CAPACITY
	if first > 0 {
		for i := 0; i < first; i++ {
			self.data[i] = nil
		}
		self.data = self.data[first:]
	}
	return nil
}

func (self *RamMetricStorage) GetMetric(tokens []common.Token, fromTime, toTime uint64) (map[string]MetricList, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	imResult := map[string]*MetricList{}
	for _, tok := range tokens {
		imResult[tok.ID] = &MetricList{}
	}
	for i := len(self.data) - 1; i >= 0; i-- {
		data := self.data[i]
		if fromTime <= data.Timestamp && data.Timestamp <= toTime {
			log.Printf("iterate over %d", data.Timestamp)
			for tok, metric := range data.Data {
				metricList, found := imResult[tok]
				if found {
					*metricList = append(*metricList, TokenMetricResponse{
						Timestamp: data.Timestamp,
						AfpMid:    metric.AfpMid,
						Spread:    metric.Spread,
					})
					log.Printf("token: %s, metricList: %+v", tok, metricList)
				}
			}
			log.Printf("result: %+v", imResult)
		} else if data.Timestamp <= fromTime {
			break
		}
	}
	result := map[string]MetricList{}
	for k, v := range imResult {
		result[k] = *v
	}
	return result, nil
}

func (self *RamMetricStorage) StorePendingTargetQty(data, dataType string) error {
	// TODO: add support ram storage
	return nil
}

func (self *RamMetricStorage) RemovePendingTargetQty() error {
	// TODO: add support ram storage
	return nil
}

func (self *RamMetricStorage) GetPendingTargetQty() (TokenTargetQty, error) {
	// TODO: add support ram storage
	result := TokenTargetQty{}
	return result, nil
}

func (self *RamMetricStorage) StoreTokenTargetQty(id, data string) error {
	self.mu.Lock()
	defer self.mu.Unlock()
	// self.tokenTargetQty = data
	return nil
}

func (self *RamMetricStorage) GetTokenTargetQty() (TokenTargetQty, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.tokenTargetQty, nil
}
