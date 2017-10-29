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
}

func NewReserveCore(blockchain Blockchain) *ReserveCore {
	return &ReserveCore{
		blockchain,
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
