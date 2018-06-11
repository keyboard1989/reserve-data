package exchange

import (
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	MustCreateTokenPair(base, quote string) common.TokenPair
	GetInternalTokens() ([]common.Token, error)
	GetAllTokens() ([]common.Token, error)
	GetTokenByID(tokenID string) (common.Token, error)
	GetFee(ex settings.ExchangeName) (common.ExchangeFees, error)
	GetMinDeposit(ex settings.ExchangeName) (common.ExchangesMinDeposit, error)
	GetDepositAddress(ex settings.ExchangeName) (common.ExchangeAddresses, error)
	GetTokenPairs(ex settings.ExchangeName) ([]common.TokenPair, error)
	GetAddress(name settings.AddressName) (ethereum.Address, error)
	UpdateDepositAddress(name settings.ExchangeName, token common.Token) error
}