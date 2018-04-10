package stat

import (
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
)

type StatStorageTest struct {
	storage StatStorage
}

const (
	TESTWALLETADDR string = "0xdd61803d4a56c597e0fc864f7a20ec7158c6cba5"
	TESTCOUNTRY    string = "unknown"
	TESTASSETADDR  string = "0x2aab2b157a03915c8a73adae735d0cf51c872f31"
)

func NewStatStorageTest(storage StatStorage) *StatStorageTest {
	return &StatStorageTest{storage}
}

func (self *StatStorageTest) TestTradeStatsSummary() error {
	var err error
	mtStat := common.MetricStats{
		10.0,
		4567.8,
		11.1,
		3,
		2,
		4,
		5,
		0,
		0,
	}
	tzmtStat := common.MetricStatsTimeZone{0: {0: mtStat}}
	updates := map[string]common.MetricStatsTimeZone{"trade_summary": tzmtStat}
	if err := self.storage.SetTradeSummary(updates); err != nil {
		return err
	}
	tradeSum, err := self.storage.GetTradeSummary(0, 86400000, 0)
	if err != nil {
		return err
	}
	//log.Println(tradeSum)
	result, ok := (tradeSum[0]).(common.MetricStats)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return wrong type")
	}
	usdVol := (result.USDVolume)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return missing field")
	}
	if usdVol != 4567.8 {
		return errors.New("Wrong usd volume value returned: %v expected 4567.8 ")

	}
	return nil
}

func (self *StatStorageTest) TestWalletStats() error {
	var err error
	mtStat := common.MetricStats{
		10.0,
		4567.8,
		11.1,
		3,
		2,
		4,
		5,
		0,
		0,
	}
	tzmtStat := common.MetricStatsTimeZone{0: {0: mtStat}}
	updates := map[string]common.MetricStatsTimeZone{TESTWALLETADDR: tzmtStat}
	err = self.storage.SetWalletStat(updates)
	if err != nil {
		return err
	}
	walletStat, err := self.storage.GetWalletStats(0, 86400000, TESTWALLETADDR, 0)
	result, ok := (walletStat[0]).(common.MetricStats)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return wrong type")
	}
	usdVol := (result.USDVolume)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return missing field")
	}
	if usdVol != 4567.8 {
		return errors.New("Wrong usd volume value returned: %v expected 4567.8 ")

	}
	return nil
}

func (self *StatStorageTest) TestCountryStats() error {
	var err error
	mtStat := common.MetricStats{
		10.0,
		4567.8,
		11.1,
		3,
		2,
		4,
		5,
		0,
		0,
	}
	tzmtStat := common.MetricStatsTimeZone{0: {0: mtStat}}
	updates := map[string]common.MetricStatsTimeZone{TESTCOUNTRY: tzmtStat}
	err = self.storage.SetWalletStat(updates)
	if err != nil {
		return err
	}
	walletStat, err := self.storage.GetCountryStats(0, 86400000, TESTCOUNTRY, 0)
	result, ok := (walletStat[0]).(common.MetricStats)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return wrong type")
	}
	usdVol := (result.USDVolume)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return missing field")
	}
	if usdVol != 4567.8 {
		return errors.New("Wrong usd volume value returned: %v expected 4567.8 ")

	}
	return nil
}

func (self *StatStorageTest) TestVolumeStats() error {
	var err error
	vlStat := common.VolumeStats{
		10.0,
		4567.8,
		11.1,
	}
	tzvlStat := common.VolumeStatsTimeZone{"D": {0: vlStat}}
	updates := map[string]common.VolumeStatsTimeZone{TESTASSETADDR: tzvlStat}
	err = self.storage.SetVolumeStat(updates)
	if err != nil {
		return err
	}
	assetVol, err := self.storage.GetAssetVolume(0, 86400000, "D", TESTASSETADDR)
	result, ok := (assetVol[0]).(common.VolumeStats)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return wrong type")
	}
	usdVol := (result.USDAmount)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return missing field")
	}
	if usdVol != 4567.8 {
		return errors.New("Wrong usd volume value returned: %v expected 4567.8 ")

	}
	assetVol, err = self.storage.GetUserVolume(0, 86400000, "D", TESTASSETADDR)
	result, ok = (assetVol[0]).(common.VolumeStats)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return wrong type")
	}
	usdVol = (result.USDAmount)
	if !ok {
		return errors.New("Type mismatched: get trade stat summary return missing field")
	}
	if usdVol != 4567.8 {
		return errors.New("Wrong usd volume value returned: %v expected 4567.8 ")

	}

	return nil

}

func (self *StatStorageTest) TestBurnFee() error {
	var err error
	return err
}

func (self *StatStorageTest) TestWalletAddress() error {
	var err error
	err = self.storage.SetWalletAddress("0xdd61803d4a56c597e0fc864f7a20ec7158c6cba5")
	if err != nil {
		return err
	}
	walletaddresses, err := self.storage.GetWalletAddress()
	if err != nil {
		return err
	}
	if len(walletaddresses) != 1 {
		return errors.New(fmt.Sprintf("expected 1 record, got %d record of wallet addresses returned", len(walletaddresses)))
	}
	if walletaddresses[0] != "0xdd61803d4a56c597e0fc864f7a20ec7158c6cba5" {
		return errors.New(fmt.Sprintf("expected address 0xdd61803d4a56c597e0fc864f7a20ec7158c6cba5, got %s instead", walletaddresses[0]))
	}
	return err
}
