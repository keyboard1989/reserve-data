package blockchain

import (
	"context"
	"math/big"

	"github.com/KyberNetwork/reserve-data/common/blockchain"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (self *Blockchain) SetBaseRate(context context.Context, opts blockchain.TxOpts, tokens []ethereum.Address, baseBuy []*big.Int, baseSell []*big.Int, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return self.BuildTx(context, opts, self.pricing, "setBaseRate", tokens, baseBuy, baseSell, buy, sell, blockNumber, indices)
}

func (self *Blockchain) SetCompactData(context context.Context, opts blockchain.TxOpts, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return self.BuildTx(context, opts, self.pricing, "setCompactData", buy, sell, blockNumber, indices)
}

func (self *Blockchain) SetImbalanceStepFunction(context context.Context, opts blockchain.TxOpts, token ethereum.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return self.BuildTx(context, opts, self.pricing, "setImbalanceStepFunction", token, xBuy, yBuy, xSell, ySell)
}

func (self *Blockchain) SetQtyStepFunction(context context.Context, opts TxOpts, token ethereum.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return self.BuildTx(context, opts, self.pricing, "setQtyStepFunction", token, xBuy, yBuy, xSell, ySell)
}

func (self *Blockchain) GetRate(context context.Context, opts blockchain.CallOpts, token ethereum.Address, currentBlockNumber *big.Int, buy bool, qty *big.Int) (*big.Int, error) {
	out := big.NewInt(0)
	err := self.Call(context, opts, self.pricing, out, "getRate", token, currentBlockNumber, buy, qty)
	return out, err
}
