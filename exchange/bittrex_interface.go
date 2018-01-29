package exchange

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type BittrexInterface interface {
	FetchOnePairData(
		pair common.TokenPair, timepoint uint64) (Bittresp, error)

	GetInfo(timepoint uint64) (Bittinfo, error)

	GetExchangeInfo() (BittExchangeInfo, error)

	GetDepositAddress(currency string) (BittrexDepositAddress, error)

	Withdraw(
		token common.Token,
		amount *big.Int,
		address ethereum.Address,
		timepoint uint64) (Bittwithdraw, error)

	Trade(
		tradeType string,
		base, quote common.Token,
		rate, amount float64,
		timepoint uint64) (Bitttrade, error)

	CancelOrder(uuid string, timepoint uint64) (Bittcancelorder, error)

	DepositHistory(currency string, timepoint uint64) (Bittdeposithistory, error)

	WithdrawHistory(currency string, timepoint uint64) (Bittwithdrawhistory, error)

	OrderStatus(uuid string, timepoint uint64) (Bitttraderesult, error)
}
