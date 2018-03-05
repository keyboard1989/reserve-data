package intermediator

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Blockchain interface {
	FetchBalanceData(addr ethereum.Address, atBlock *big.Int, timepoint uint64) (map[string]common.BalanceEntry, error)
	SendFromAccountToExchange(amount *big.Int, exchangeAddress ethereum.Address) (*types.Transaction, error)
}
