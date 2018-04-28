package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type GlobalStorage interface {
	StoreGoldInfo(data common.GoldData) error
}
