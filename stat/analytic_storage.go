package stat

type AnalyticStorage interface {
	UpdatePriceAnalyticData(timestamp uint64, value []byte) error
	GetPriceAnalyticData(fromTime uint64, toTime uint64) (map[string][]byte, error)
}
