package bittrex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
)

func TestBittrexStorage(t *testing.T) {
	boltFile := "test_bittrex_storage.db"
	tmpDir, err := ioutil.TempDir("", "bittrex_storage")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}
	}()

	exchangeTradeHistory := common.ExchangeTradeHistory{
		common.TokenPairID("OMGETH"): []common.TradeHistory{
			{
				ID:        "12342",
				Price:     0.132131,
				Qty:       12.3123,
				Type:      "buy",
				Timestamp: 1528949872,
			},
		},
	}

	storage, err := NewBoltStorage(filepath.Join(tmpDir, boltFile))
	if err != nil {
		t.Fatalf("Could not init bittrex bolt storage: %s", err.Error())
	}

	// store trade history
	err = storage.StoreTradeHistory(exchangeTradeHistory)
	if err != nil {
		t.Fatal(err)
	}

	// get trade history
	var tradeHistory common.ExchangeTradeHistory
	fromTime := uint64(1528934400000)
	toTime := uint64(1529020800000)
	tradeHistory, err = storage.GetTradeHistory(fromTime, toTime)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tradeHistory, exchangeTradeHistory) {
		t.Fatal("Bittrex get wrong trade history")
	}

	// get last trade history id
	var lastHistoryID string
	lastHistoryID, err = storage.GetLastIDTradeHistory("OMGETH")
	if err != nil {
		t.Fatalf("Bittrex get last trade history id error: %s", err.Error())
	}
	if lastHistoryID != "12342" {
		t.Fatalf("Bittrex get last trade history wrong, shoudl return 12342 receive: %s", lastHistoryID)
	}

	newDepositActivityID := common.NewActivityID(
		1513328774800747341,
		"0x4e3c6c5e5c56ef2f65867e0dac874f3e2ec2f66e05793b7a52281549c02e68d9|OMG|5",
	)
	// test new deposit
	checkNewDeposit := storage.IsNewBittrexDeposit(newDepositActivityID.Timepoint, newDepositActivityID)
	if !checkNewDeposit {
		t.Fatal("Bittrex check new deposit wrong.")
	}

	// register new deposit
	err = storage.RegisterBittrexDeposit(newDepositActivityID.Timepoint, newDepositActivityID)
	if err != nil {
		t.Fatalf("Bittrex register new deposit wrong: %s", err.Error())
	}

	// checkNewDeposit = storage.IsNewBittrexDeposit(newDepositActivityID.Timepoint, newDepositActivityID)
	// if checkNewDeposit {
	// 	t.Fatal("Bittrex check new deposit wrong")
	// }
}
