package huobi

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Storage interface {
	StoreIntermediateTx(id common.ActivityID, data common.TXEntry) error
	StorePendingIntermediateTx(id common.ActivityID, data common.TXEntry) error

	GetIntermedatorTx(id common.ActivityID) (common.TXEntry, error)
	GetPendingIntermediateTXs() (map[common.ActivityID]common.TXEntry, error)
	RemovePendingIntermediateTx(id common.ActivityID) error
}
