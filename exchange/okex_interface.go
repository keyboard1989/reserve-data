package exchange

import "github.com/KyberNetwork/reserve-data/common"

type OkexInterface interface {
	GetDepthOnePair(
		pair common.TokenPair, timepoint uint64) (OkexDepth, error)

	GetExchangeInfo() (OkexInfo, error)
}
