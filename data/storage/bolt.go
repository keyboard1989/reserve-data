package storage

const (
	PRICE_BUCKET            string = "prices"
	BALANCE_BUCKET          string = "balances"
	EXCHANGE_BALANCE_BUCKET string = "ebalances"
	RATE_BUCKET             string = "rates"
)

type BoltStorage struct {
}

func (self *BoltStorage) CurrentPriceVersion(timepoint uint64) (common.Version, error) {
}

func (self *BoltStorage) GetAllPrices(common.Version) (map[common.TokenPairID]common.OnePrice, error) {
}

func (self *BoltStorage) GetOnePrice(common.TokenPairID, common.Version) (common.OnePrice, error) {
}

func (self *BoltStorage) CurrentBalanceVersion(timepoint uint64) (common.Version, error) {
}
func (self *BoltStorage) GetAllBalances(common.Version) (map[string]common.BalanceEntry, error) {
}

func (self *BoltStorage) CurrentEBalanceVersion(timepoint uint64) (common.Version, error) {
}
func (self *BoltStorage) GetAllEBalances(common.Version) (map[common.ExchangeID]common.EBalanceEntry, error) {
}

func (self *BoltStorage) CurrentRateVersion(timepoint uint64) (common.Version, error) {
}
func (self *BoltStorage) GetAllRates(common.Version) (common.AllRateEntry, error) {
}

func (self *BoltStorage) StorePrice(data map[common.TokenPairID]common.OnePrice, timepoint uint64) error {
}
func (self *BoltStorage) StoreBalance(data map[string]common.BalanceEntry, timepoint uint64) error {
}
func (self *BoltStorage) StoreEBalance(data map[common.ExchangeID]common.EBalanceEntry, timepoint uint64) error {
}
func (self *BoltStorage) StoreRate(data common.AllRateEntry, timepoint uint64) error {
}
