package core

import (
	"math/big"
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type testExchange struct {
}

func (self testExchange) ID() common.ExchangeID {
	return "bittrex"
}
func (self testExchange) Address(token common.Token) (address ethereum.Address, supported bool) {
	return ethereum.Address{}, true
}
func (self testExchange) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	return "withdrawid", nil
}
func (self testExchange) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	return "tradeid", 10, 5, false, nil
}
func (self testExchange) CancelOrder(id common.ActivityID) error {
	return nil
}
func (self testExchange) MarshalText() (text []byte, err error) {
	return []byte("bittrex"), nil
}

type testBlockchain struct {
}

func (self testBlockchain) Send(
	token common.Token,
	amount *big.Int,
	address ethereum.Address) (ethereum.Hash, error) {
	return ethereum.Hash{}, nil
}

func (self testBlockchain) SetRates(
	sources []ethereum.Address,
	dests []ethereum.Address,
	rates []*big.Int,
	expiryBlocks []*big.Int) (ethereum.Hash, error) {
	return ethereum.Hash{}, nil
}

type testActivityStorage struct {
	PendingDeposit bool
}

func (self testActivityStorage) Record(
	action string,
	id common.ActivityID,
	destination string,
	params map[string]interface{},
	result map[string]interface{},
	estatus string,
	mstatus string,
	timepoint uint64) error {
	return nil
}

func (self testActivityStorage) HasPendingDeposit(
	token common.Token, exchange common.Exchange) bool {
	if token.ID == "OMG" && exchange.ID() == "bittrex" {
		return self.PendingDeposit
	} else {
		return false
	}
}

func getTestCore(hasPendingDeposit bool) *ReserveCore {
	return NewReserveCore(
		testBlockchain{},
		testActivityStorage{hasPendingDeposit},
		ethereum.Address{},
	)
}

func TestNotAllowDeposit(t *testing.T) {
	alreadyHasDepositForOMGOnBittrex := true
	core := getTestCore(alreadyHasDepositForOMGOnBittrex)
	_, err := core.Deposit(
		testExchange{},
		common.Token{"OMG", "0x1111111111111111111111111111111111111111", 18},
		big.NewInt(10),
		common.GetTimepoint(),
	)
	if err == nil {
		t.Fatalf("Expected to return an error protecting user from deposit when there is another pending deposit")
	}
	_, err = core.Deposit(
		testExchange{},
		common.Token{"KNC", "0x1111111111111111111111111111111111111111", 18},
		big.NewInt(10),
		common.GetTimepoint(),
	)
	if err != nil {
		t.Fatalf("Expected to be able to deposit different token")
	}
}
