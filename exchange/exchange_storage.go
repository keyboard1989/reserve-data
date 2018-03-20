package exchange

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StoreIntermediateTx(hash string, exchangeID string, tokenID string, miningStatus string, exchangeStatus string, Amount float64, Timestamp common.Timestamp, id common.ActivityID) error
	StorePendingIntermediateTx(hash string, exchangeID string, tokenID string, miningStatus string, exchangeStatus string, Amount float64, Timestamp common.Timestamp, id common.ActivityID) error

	GetIntermedatorTx(id common.ActivityID) (common.TXEntry, error)
	GetPendingIntermediateTXs(TimeStamp uint64) (map[common.ActivityID]common.TXEntry, error)
	RemovePendingIntermediateTx(id common.ActivityID) error
}
