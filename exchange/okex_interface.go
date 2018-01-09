package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type OkexInterface interface {
	GetDepthOnePair(
		pair common.TokenPair, timepoint uint64) (OkexDepth, error)

	GetExchangeInfo() (OkexInfo, error)

	GetInfo(timepoint uint64) (OkexAccountInfo, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address,
		timepoint uint64) (string, error)

	Trade(base, quote common.Token, rate, amount float64, timepoint uint64) (OkexTrade, error)

	CancelOrder(symbol string, id uint64, timepoint uint64) (OkexCancel, error)

	OrderStatus(symbol string, id uint64, timepoint uint64) (OkexOrderStatus, error)
}
