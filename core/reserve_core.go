package core

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ReserveCore struct {
	blockchain Blockchain
	rm         ethereum.Address
}

func NewReserveCore(blockchain Blockchain, rm ethereum.Address) *ReserveCore {
	return &ReserveCore{
		blockchain,
		rm,
	}
}

func (self ReserveCore) Deposit(
	exchange common.Exchange,
	token common.Token,
	amount *big.Int) (ethereum.Hash, error) {

	address, supported := exchange.Address(token)
	if !supported {
		return ethereum.Hash{}, errors.New(fmt.Sprintf("Exchange %s doesn't support token %s", exchange.ID(), token.ID))
	}
	return self.blockchain.Send(token, amount, address)
}

func (self ReserveCore) Withdraw(
	exchange common.Exchange, token common.Token,
	amount *big.Int) error {

	_, supported := exchange.Address(token)
	if !supported {
		return errors.New(fmt.Sprintf("Exchange %s doesn't support token %s", exchange.ID(), token.ID))
	}
	return exchange.Withdraw(token, amount, self.rm)
}

func (self ReserveCore) SetRates(
	sources []common.Token,
	dests []common.Token,
	rates []*big.Int,
	expiryBlocks []*big.Int) (ethereum.Hash, error) {

	lensources := len(sources)
	lendests := len(dests)
	lenrates := len(rates)
	lenblocks := len(expiryBlocks)
	if lensources != lendests || lensources != lenrates || lensources != lenblocks {
		return ethereum.Hash{}, errors.New("Sources, dests, rates and expiryBlocks must have the same length")
	} else {
		sourceAddrs := []ethereum.Address{}
		for _, source := range sources {
			sourceAddrs = append(sourceAddrs, ethereum.HexToAddress(source.Address))
		}
		destAddrs := []ethereum.Address{}
		for _, dest := range dests {
			destAddrs = append(destAddrs, ethereum.HexToAddress(dest.Address))
		}
		return self.blockchain.SetRates(sourceAddrs, destAddrs, rates, expiryBlocks)
	}
}
