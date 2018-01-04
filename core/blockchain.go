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
		tokens []ethereum.Address,
		buys []*big.Int,
		sells []*big.Int,
		block *big.Int) (ethereum.Hash, error)
}
