package exchange

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type testBittrexInterface struct{}

func (self testBittrexInterface) FetchOnePairData(wq *sync.WaitGroup, pair common.TokenPair, data *sync.Map, timepoint uint64) {
}
func (self testBittrexInterface) GetInfo(timepoint uint64) (Bittinfo, error) {
	return Bittinfo{}, nil
}
func (self testBittrexInterface) Withdraw(
	token common.Token,
	amount *big.Int,
	address ethereum.Address,
	timepoint uint64) (Bittwithdraw, error) {
	return Bittwithdraw{}, nil
}
func (self testBittrexInterface) Trade(
	tradeType string,
	base, quote common.Token,
	rate, amount float64,
	timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	return "", 0, 0, false, nil
}
func (self testBittrexInterface) CancelOrder(uuid string, timepoint uint64) (Bittcancelorder, error) {
	return Bittcancelorder{}, nil
}
func (self testBittrexInterface) DepositHistory(currency string, timepoint uint64) (Bittdeposithistory, error) {
	str := `{"success":true,"message":"","result":[{"Id":46452872,"Amount":5.00000000,"Currency":"OMG","Confirmations":42,"LastUpdated":"2017-12-15T09:27:09.597","TxId":"0x15ccaab008f161efeee0febc3e32242846cea1fc93995e5abc6fb88d94ae7d21","CryptoAddress":"0x9db6e8d2d133448dbcf755f19d540253da4ba043"}]}`
	res := Bittdeposithistory{}
	err := json.Unmarshal([]byte(str), &res)
	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", res)
	return res, err
}
func (self testBittrexInterface) WithdrawHistory(currency string, timepoint uint64) (Bittwithdrawhistory, error) {
	return Bittwithdrawhistory{}, nil
}
func (self testBittrexInterface) OrderStatus(uuid string, timepoint uint64) (Bitttraderesult, error) {
	return Bitttraderesult{}, nil
}

type testBittrexStorage struct{}

func (self testBittrexStorage) IsNewBittrexDeposit(id uint64) bool {
	return true
}

func (self testBittrexStorage) RegisterBittrexDeposit(id uint64) error {
	return nil
}

func getTestBittrex() *Bittrex {
	return &Bittrex{
		testBittrexInterface{},
		[]common.TokenPair{},
		map[string]ethereum.Address{},
		testBittrexStorage{},
	}
}

func TestDepositStatus(t *testing.T) {
	activityID := common.ActivityID{
		1513328774800747341,
		"0x4e3c6c5e5c56ef2f65867e0dac874f3e2ec2f66e05793b7a52281549c02e68d9|OMG|5",
	}
	bitt := getTestBittrex()
	out, err := bitt.DepositStatus(activityID, common.GetTimepoint())
	if err != nil {
		t.Fatalf("Expected convert successfully but got error: %v", err)
	} else {
		if out != "done" {
			t.Fatalf("Expected done, got %v", out)
		}
	}
}
