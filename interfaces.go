package reserve

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// all of the functions must support concurrency
type ReserveData interface {
	CurrentPriceVersion() (common.Version, error)
	GetAllPrices() (common.AllPriceResponse, error)
	GetOnePrice(common.TokenPairID) (common.OnePriceResponse, error)

	CurrentBalanceVersion() (common.Version, error)
	GetAllBalances() (common.AllBalanceResponse, error)
	Run() error
}
