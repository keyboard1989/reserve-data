package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Blockchain interface {
	FetchBalanceData(addr ethereum.Address, timepoint uint64) (map[string]common.BalanceEntry, error)
	FetchRates(sources, dests []common.Token, timepoint uint64) (common.AllRateEntry, error)
}
