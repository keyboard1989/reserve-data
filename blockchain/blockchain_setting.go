package blockchain

import (
	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Setting interface {
	GetInternalTokenByID(tokenID string) (common.Token, error)
	GetInternalTokens() ([]common.Token, error)
	ETHToken() common.Token
	AddAddressToSet(setName string, address ethereum.Address) error
	GetAddress(name string) (ethereum.Address, error)
	GetAddresses(setName string) ([]ethereum.Address, error)
}
