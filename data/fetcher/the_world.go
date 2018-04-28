package fetcher

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type TheWorld interface {
	GetGoldInfo() (common.GoldData, error)
}
