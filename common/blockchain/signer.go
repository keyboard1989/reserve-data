package blockchain

import (
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Signer interface {
	GetAddress() ethereum.Address
	Sign(*types.Transaction) (*types.Transaction, error)
}
