package blockchain

import (
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

func (self *Blockchain) GetBalances(context context.Context, opts CallOpts, reserve ethereum.Address, tokens []ethereum.Address) ([]*big.Int, error) {
	out := new([]*big.Int)
	err := self.Call(context, opts, self.wrapper, out, "getBalances", reserve, tokens)
	return *out, err
}

func (self *Blockchain) GetTokenIndicies(context context.Context, opts CallOpts, atBlock *big.Int, ratesContract ethereum.Address, tokenList []ethereum.Address) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := self.Call(context, opts, self.wrapper, out, "getTokenIndicies", ratesContract, tokenList)
	return *ret0, *ret1, err
}

func (self *Blockchain) GetTokenRates(
	context context.Context, opts CallOpts,
	ratesContract ethereum.Address,
	tokenList []ethereum.Address) ([]*big.Int, []*big.Int, []int8, []int8, []*big.Int, error) {
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
	err := self.Call(context, opts, self.wrapper, out, "getTokenRates", ratesContract, tokenList)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

func (self *Blockchain) GetReserveRates(
	context context.Context, opts CallOpts,
	reserveAddress ethereum.Address,
	srcAddresses []ethereum.Address,
	destAddresses []ethereum.Address) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := self.Call(context, opts, self.wrapper, out, "getReserveRate", reserveAddress, srcAddresses, destAddresses)
	if err != nil {
		log.Println("cannot get reserve rates: ", err.Error())
	}
	return *ret0, *ret1, err
}
