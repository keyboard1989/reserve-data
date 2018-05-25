package testutil

import (
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
)

// GlobalStorageTestSuite is a test helper for data.GlobalStorage and fetcher.GlobalStorage.
type GlobalStorageTestSuite struct {
	t   *testing.T
	dgs data.GlobalStorage
	fgs fetcher.GlobalStorage
}

// NewGlobalStorageTestSuite creates new a new test suite with given implementations of data.GlobalStorage and
// fetcher.GlobalStorage.
func NewGlobalStorageTestSuite(t *testing.T, dgs data.GlobalStorage, fgs fetcher.GlobalStorage) *GlobalStorageTestSuite {
	t.Helper()
	return &GlobalStorageTestSuite{
		t:   t,
		dgs: dgs,
		fgs: fgs,
	}
}

// Run executes test suite.
func (ts *GlobalStorageTestSuite) Run() {
	ts.t.Helper()

	dgxSts := []string{"first", "second", "third"}
	for _, dgxSt := range dgxSts {
		err := ts.fgs.StoreGoldInfo(common.GoldData{
			Timestamp: common.GetTimepoint(),
			DGX: common.DGXGoldData{
				Valid:  true,
				Status: dgxSt,
			},
		})
		if err != nil {
			ts.t.Fatal(err)
		}
	}

	ts.t.Log("getting gold info of future timepoint, expected to get latest version")
	var futureTimepoint uint64 = common.GetTimepoint() + 100
	version, err := ts.dgs.CurrentGoldInfoVersion(futureTimepoint)
	if err != nil {
		ts.t.Fatal(err)
	}

	goldInfo, err := ts.dgs.GetGoldInfo(version)
	if err != nil {
		ts.t.Fatal(err)
	}
	ts.t.Logf("latest gold info version: %d", version)

	lastDgxSt := dgxSts[(len(dgxSts) - 1)]
	if goldInfo.DGX.Status != lastDgxSt {
		ts.t.Errorf("getting wrong dgx status, expected: %s, got: %s", lastDgxSt, goldInfo.DGX.Status)
	}
}
