package reserve

import (
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"math/big"
)

// all of the functions must support concurrency
type ReserveData interface {
	CurrentPriceVersion() (common.Version, error)
	GetAllPrices() (common.AllPriceResponse, error)
	GetOnePrice(common.TokenPairID) (common.OnePriceResponse, error)

	CurrentBalanceVersion() (common.Version, error)
	GetAllBalances() (common.AllBalanceResponse, error)

	CurrentEBalanceVersion() (common.Version, error)
	GetAllEBalances() (common.AllEBalanceResponse, error)

	CurrentRateVersion() (common.Version, error)
	GetAllRates() (common.AllRateResponse, error)
	Run() error
}

type ReserveCore interface {
	// place order
	// cancel order
	Deposit(exchange common.Exchange, token common.Token, amount *big.Int) (ethereum.Hash, error)
	// withdraw
	Withdraw(exchange common.Exchange, token common.Token, amount *big.Int) error

	// blockchain related action
	SetRates(sources []common.Token, dests []common.Token, rates []*big.Int, expiryBlocks []*big.Int) (ethereum.Hash, error)
}
