package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type BinanceInterface interface {
	GetDepthOnePair(
		pair common.TokenPair, timepoint uint64) (Binaresp, error)

	OpenOrdersForOnePair(
		pair common.TokenPair, timepoint uint64) (Binaorders, error)

	GetInfo(timepoint uint64) (Binainfo, error)

	GetExchangeInfo() (BinanceExchangeInfo, error)

	GetDepositAddress(tokenID string) (Binadepositaddress, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address,
		timepoint uint64) (string, error)

	Trade(
		tradeType string,
		base, quote common.Token,
		rate, amount float64,
		timepoint uint64) (Binatrade, error)

	CancelOrder(symbol string, id uint64) (Binacancel, error)

	DepositHistory(startTime, endTime uint64) (Binadeposits, error)

	WithdrawHistory(
		startTime, endTime uint64) (Binawithdrawals, error)

	OrderStatus(
		symbol string, id uint64, timepoint uint64) (Binaorder, error)

	GetServerTime() (uint64, error)
}
