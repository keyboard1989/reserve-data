package blockchain

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetInternalTokens() ([]common.Token, error)
	ETHToken() common.Token
}
