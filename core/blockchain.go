package core

import (
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Blockchain interface {
	Send(
		token common.Token,
		amount *big.Int,
		address ethereum.Address) (ethereum.Hash, error)
	SetRates(
		sources []ethereum.Address,
		dests []ethereum.Address,
		rates []*big.Int,
		expiryBlocks []*big.Int) (ethereum.Hash, error)
}
