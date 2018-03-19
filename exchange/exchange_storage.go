package exchange

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StoreIntermediateTx(hash string, exchangeID string, tokenID string, miningStatus string, exchangeStatus string, Amount float64, Timestamp common.Timestamp, id common.ActivityID) error
	GetIntermedatorTx(id common.ActivityID) (common.TXEntry, error)
	GetPendingIntermediateTXID(TimeStamp uint64) ([]common.ActivityID, error)
	StorePendingIntermediateTxID(Timestamp uint64, id common.ActivityID) error
}
