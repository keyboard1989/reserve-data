package data

import "github.com/KyberNetwork/reserve-data/common"

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
}
