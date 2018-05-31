package stat

import (
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Setting interface {
	GetInternalTokens() ([]common.Token, error)
	GetActiveTokens() ([]common.Token, error)
	GetTokenByAddress(addr ethereum.Address) (common.Token, error)
	ETHToken() common.Token
	GetActiveTokenByID(id string) (common.Token, error)
}