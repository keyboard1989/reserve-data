package configuration

import (
	"os"
	"testing"

	"github.com/KyberNetwork/reserve-data/stat"
	statstorage "github.com/KyberNetwork/reserve-data/stat/storage"
)

func SetupBoltStorageTester(name string) (*stat.StorageTest, error) {
	storage, err := statstorage.NewBoltStorage(
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/configuration/" + name,
	)
	if err != nil {
		return nil, err
	}
	return stat.NewStorageTest(storage), nil
}

func TearDownBolt(name string) {
	os.Remove(
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/configuration/" + name,
	)
}

func TestBoltAsStatStorage(t *testing.T) {
	dbname := "test1.db"
	tester, err := SetupBoltStorageTester(dbname)
	if err != nil {
		t.Fatalf("Testing bolt as a stat storage: init failed(%s)", err)
	}
	defer TearDownBolt(dbname)
	if err = tester.TestStoreCatLog(); err != nil {
		t.Fatalf("Testing bolt as a stat storage: test store cat log failed(%s)", err)
	}
	if err = tester.TestStoreCatLogThenUpdateUserAddresses(); err != nil {
		t.Fatalf("Testing bolt as a stat storage: test store cat log and then update user addresses failed(%s)", err)
	}
	if err = tester.TestUpdateUserAddressesThenStoreCatLog(); err != nil {
		t.Fatalf("Testing bolt as a stat storage: test update user addresses and then store cat log failed(%s)", err)
	}
}
