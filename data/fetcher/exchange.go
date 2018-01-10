package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Exchange interface {
	ID() common.ExchangeID
	Name() string
	TokenPairs() []common.TokenPair
	FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error)
	FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error)

	// OrderStatus checks status of an order corresponding to activity id.
	// returns: status (string) and error.
	//   - error: is not nil if there is at least one error communicating to the exchange, error is nil otherwise
	//   - status: can be "" or "done". "" means the order is still waiting to be fulfilled. "done" means the order
	// is fulfilled by the exchange or canceled by us
	OrderStatus(id common.ActivityID, timepoint uint64) (string, error)

	// DepositStatus check status of a deposit
	// returns: status (string) and error.
	//   - error: is not nil if there is at least one error communicating to the exchange, error is nil otherwise
	//   - status: can be "" or "done". "" means the deposit is pending. "done" means the deposit is done.
	DepositStatus(id common.ActivityID, timepoint uint64) (string, error)

	// WithdraStatus check status of a withdraw
	// returns: status (string), txid (string) and error.
	//   - error: is not nil if there is at least one error communicating to the exchange, error is nil otherwise
	//   - status: can be "" or "done". "" means the withdraw is pending. "done" means the withdraw is done.
	//   - txid: is the tx id that the exchange uses to send back the funds
	WithdrawStatus(id common.ActivityID, timepoint uint64) (string, string, error)
}
