package http

import (
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetActiveTokenByID(tokenID string) (common.Token, error)
	GetInternalTokens() ([]common.Token, error)
	GetAllTokens() ([]common.Token, error)
	NewTokenPairFromID(base, quote string) (common.TokenPair, error)
	UpdateToken(t common.Token) error
	AddAddressToSet(setName settings.AddressSetName, address ethereum.Address) error
	UpdateAddress(name settings.AddressName, address ethereum.Address) error
}
