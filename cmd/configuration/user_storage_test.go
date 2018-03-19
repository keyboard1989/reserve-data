package configuration

import (
	"os"
	"testing"

	"github.com/KyberNetwork/reserve-data/stat"
	statstorage "github.com/KyberNetwork/reserve-data/stat/storage"
)

func SetupBoltUserStorageTester(name string) (*stat.UserStorageTest, error) {
	storage, err := statstorage.NewBoltUserStorage(
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/configuration/" + name,
	)
	if err != nil {
		return nil, err
	}
	return stat.NewUserStorageTest(storage), nil
}

func TearDownBolt(name string) {
	os.Remove(
		"/go/src/github.com/KyberNetwork/reserve-data/cmd/configuration/" + name,
	)
}

func doOneTest(f func(tester *stat.UserStorageTest, t *testing.T), t *testing.T) {
	dbname := "test1.db"
	tester, err := SetupBoltUserStorageTester(dbname)
	if err != nil {
		t.Fatalf("Testing bolt as a stat storage: init failed(%s)", err)
	}
	defer TearDownBolt(dbname)
	f(tester, t)
}

func TestBoltAsStatStorage(t *testing.T) {
	doOneTest(func(tester *stat.UserStorageTest, t *testing.T) {
		if err := tester.TestUpdateAddressCategory(); err != nil {
			t.Fatalf("Testing bolt as a stat storage: test update address category failed(%s)", err)
		}
	}, t)
	doOneTest(func(tester *stat.UserStorageTest, t *testing.T) {
		if err := tester.TestUpdateUserAddressesThenUpdateAddressCategory(); err != nil {
			t.Fatalf("Testing bolt as a stat storage: test update user addresses and then update address category failed(%s)", err)
		}
	}, t)
	doOneTest(func(tester *stat.UserStorageTest, t *testing.T) {
		if err := tester.TestUpdateAddressCategoryThenUpdateUserAddresses(); err != nil {
			t.Fatalf("Testing bolt as a stat storage: test update address category and then update user addresses failed(%s)", err)
		}
	}, t)
}
