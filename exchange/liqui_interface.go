package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type LiquiInterface interface {
	Depth(tokens string, timepoint uint64) (Liqresp, error)

	GetInfo(timepoint uint64) (Liqinfo, error)

	ActiveOrders(timepoint uint64) (Liqorders, error)

	OrderInfo(orderID string, timepoint uint64) (Liqorderinfo, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address,
		timepoint uint64) error

	Trade(
		tradeType string,
		base, quote common.Token,
		rate, amount float64,
		timepoint uint64) (id string, done float64, remaining float64, finished bool, err error)

	CancelOrder(id string) (Liqcancel, error)
}
