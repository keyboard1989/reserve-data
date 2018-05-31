package http

import "github.com/KyberNetwork/reserve-data/common"

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetActiveTokenByID(tokenID string) (common.Token, error)
	GetInternalTokens() ([]common.Token, error)
	GetAllTokens() ([]common.Token, error)
	NewTokenPair(base, quote string) (common.TokenPair, error)
	UpdateToken(t common.Token) error
}
