package market

var StorageInstance Storage

type Storage interface {
	GetVersion() int64
	GetTimestamp(int64) string
	GetOnePairData(string, string, int64) (OnePairData, error)
	GetAllPairData(int64) (AllPriceData, error)

	StoreNewData(timestamp string, data AllPriceData) error
}
