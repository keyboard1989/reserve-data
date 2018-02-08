package blockchain

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type KNWrapperContract struct {
	*KNContractBase
}

func (self *KNWrapperContract) GetBalances(opts *bind.CallOpts, atBlock *big.Int, reserve ethereum.Address, tokens []ethereum.Address) ([]*big.Int, error) {
	out := new([]*big.Int)
	err := self.KNContractBase.Call(opts, atBlock, out, "getBalances", reserve, tokens)
	return *out, err
}

func (self *KNWrapperContract) GetTokenIndicies(opts *bind.CallOpts, atBlock *big.Int, ratesContract ethereum.Address, tokenList []ethereum.Address) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := self.KNContractBase.Call(opts, atBlock, out, "getTokenIndicies", ratesContract, tokenList)
	return *ret0, *ret1, err
}

func (self *KNWrapperContract) GetTokenRates(opts *bind.CallOpts, atBlock *big.Int, ratesContract ethereum.Address, tokenList []ethereum.Address) ([]*big.Int, []*big.Int, []int8, []int8, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
		ret2 = new([]int8)
		ret3 = new([]int8)
		ret4 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := self.KNContractBase.Call(opts, atBlock, out, "getTokenRates", ratesContract, tokenList)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

func NewKNWrapperContract(address ethereum.Address, client *ethclient.Client) (*KNWrapperContract, error) {
	file, err := os.Open(
		"/go/src/github.com/KyberNetwork/reserve-data/blockchain/wrapper.abi")
	if err != nil {
		return nil, err
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		return nil, err
	}
	return &KNWrapperContract{
		NewKNContractBase(address, parsed, client),
	}, nil
}
