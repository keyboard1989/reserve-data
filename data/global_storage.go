package data

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type GlobalStorage interface {
	GetGoldInfo(version common.Version) (common.GoldData, error)
	CurrentGoldInfoVersion(timepoint uint64) (common.Version, error)
}
