package reserve

import (
	"github.com/KyberNetwork/reserve-data/common"
	"math/big"
)

// all of the functions must support concurrency
type ReserveData interface {
	CurrentPriceVersion(timestamp uint64) (common.Version, error)
	GetAllPrices(timestamp uint64) (common.AllPriceResponse, error)
	GetOnePrice(id common.TokenPairID, timestamp uint64) (common.OnePriceResponse, error)

	CurrentBalanceVersion(timestamp uint64) (common.Version, error)
	GetAllBalances(timestamp uint64) (common.AllBalanceResponse, error)

	CurrentEBalanceVersion(timestamp uint64) (common.Version, error)
	GetAllEBalances(timestamp uint64) (common.AllEBalanceResponse, error)

	CurrentRateVersion(timestamp uint64) (common.Version, error)
	GetAllRates(timestamp uint64) (common.AllRateResponse, error)

	GetRecords() ([]common.ActivityRecord, error)
	GetPendingActivities() ([]common.ActivityRecord, error)

	Run() error
	Stop() error
}

type ReserveCore interface {
	// place order
	Trade(
		exchange common.Exchange,
		tradeType string,
		base common.Token,
		quote common.Token,
		rate float64,
		amount float64,
		timestamp uint64) (id common.ActivityID, done float64, remaining float64, finished bool, err error)

	Deposit(
		exchange common.Exchange,
		token common.Token,
		amount *big.Int,
		timestamp uint64) (common.ActivityID, error)

	Withdraw(
		exchange common.Exchange,
		token common.Token,
		amount *big.Int,
		timestamp uint64) (common.ActivityID, error)

	CancelOrder(id common.ActivityID, exchange common.Exchange) error

	// blockchain related action
	SetRates(sources []common.Token, dests []common.Token, rates []*big.Int, expiryBlocks []*big.Int) (common.ActivityID, error)

	// action: ['deposit', 'withdraw', 'setrate', 'order']
	// id: ID of the action, for deposit and setrate, it should be tx hash, for withdraw, order, it should be id returned
	// from the exchanges
	// return string: should be in ["", "submitted", "mined", "done", "failed"]
	// ActivityStatus(action string, id string, destination string) (string, error)
}
