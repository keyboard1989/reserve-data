package blockchain

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type KNPricingContract struct {
	*KNContractBase
}

func (self *KNPricingContract) SetBaseRate(opts *bind.TransactOpts, tokens []ethereum.Address, baseBuy []*big.Int, baseSell []*big.Int, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return self.KNContractBase.BuildTx(opts, "setBaseRate", tokens, baseBuy, baseSell, buy, sell, blockNumber, indices)
}

func (self *KNPricingContract) SetCompactData(opts *bind.TransactOpts, buy [][14]byte, sell [][14]byte, blockNumber *big.Int, indices []*big.Int) (*types.Transaction, error) {
	return self.KNContractBase.BuildTx(opts, "setCompactData", buy, sell, blockNumber, indices)
}

func (self *KNPricingContract) SetImbalanceStepFunction(opts *bind.TransactOpts, token ethereum.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return self.KNContractBase.BuildTx(opts, "setImbalanceStepFunction", token, xBuy, yBuy, xSell, ySell)
}

func (self *KNPricingContract) SetQtyStepFunction(opts *bind.TransactOpts, token ethereum.Address, xBuy []*big.Int, yBuy []*big.Int, xSell []*big.Int, ySell []*big.Int) (*types.Transaction, error) {
	return self.KNContractBase.BuildTx(opts, "setQtyStepFunction", token, xBuy, yBuy, xSell, ySell)
}

func (self *KNPricingContract) GetRate(opts *bind.CallOpts, atBlock *big.Int, token ethereum.Address, currentBlockNumber *big.Int, buy bool, qty *big.Int) (*big.Int, error) {
	out := big.NewInt(0)
	err := self.KNContractBase.Call(opts, atBlock, out, "getRate", token, currentBlockNumber, buy, qty)
	return out, err
}

func NewKNPricingContract(address ethereum.Address, client *ethclient.Client) (*KNPricingContract, error) {
	file, err := os.Open(
		"/go/src/github.com/KyberNetwork/reserve-data/blockchain/pricing.abi")
	if err != nil {
		return nil, err
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		return nil, err
	}
	return &KNPricingContract{
		NewKNContractBase(address, parsed, client),
	}, nil
}
