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
	Run() error
}

type ReserveCore interface {
	// withdraw
	// place order
	// cancel order
	// deposit
	Deposit(exchange common.Exchange, token common.Token, amount *big.Int) (ethereum.Hash, error)
	SetRates(sources []common.Token, dests []common.Token, rates []*big.Int, expiryBlocks []*big.Int) (ethereum.Hash, error)
}
