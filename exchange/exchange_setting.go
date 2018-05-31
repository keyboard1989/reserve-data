package exchange

import "github.com/KyberNetwork/reserve-data/common"

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	MustCreateTokenPair(base, quote string) common.TokenPair
	GetInternalTokens() ([]common.Token, error)
	GetTokens() ([]common.Token, error)
	GetTokenByID(tokenID string) (common.Token, error)
}
