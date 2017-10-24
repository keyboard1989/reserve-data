package alpha

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	CurrentPriceVersion() (common.Version, error)
	GetAllPrices(common.Version) (map[common.TokenPairID]common.OnePrice, error)
	GetOnePrice(common.TokenPairID, common.Version) (common.OnePrice, error)

	CurrentBalanceVersion() (common.Version, error)
	GetAllBalances(common.Version) (map[string]common.BalanceEntry, error)
}

type Fetcher interface {
	Run() error
}
