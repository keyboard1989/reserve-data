package blockchain

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (self *Blockchain) Withdraw(context context.Context, opts TxOpts, token ethereum.Address, amount *big.Int, destination ethereum.Address) (*types.Transaction, error) {
	return self.BuildTx(context, opts, self.reserve, "withdraw", token, amount, destination)
}
