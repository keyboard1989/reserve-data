package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Blockchain interface {
	FetchBalanceData(ethereum.Address) (map[string]common.BalanceEntry, error)
	FetchRates(sources, dests []common.Token) (map[common.TokenPairID]common.RateEntry, error)
}
