package stat

import (
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// This test type enforces necessary logic required for a stat storage.
// - It requires an actual storage instance to be able to run the tests.
// - It DOESNT do any tear up or tear down processes.
// - Each of its functions is for one test and will return non-nil error
// if the test didn't pass.
// - It is supposed to be used in a package that has the knowledge
// of actual storage being used as this interface.
// Eg. It should be used in cmd package where we decide to use
// bolt (for example) as the storage for stat storage
type StorageTest struct {
	storage Storage
}

func NewStorageTest(storage Storage) *StorageTest {
	return &StorageTest{storage}
}

func NewCatLog(address string, cat string) common.SetCatLog {
	return common.SetCatLog{
		0, 0,
		ethereum.HexToAddress(address),
		cat,
	}
}

func (self *StorageTest) TestStoreCatLog() error {
	lowercaseAddr := "0x8180a5ca4e3b94045e05a9313777955f7518d757"
	lowercaseCat := "0x4a"
	addr := "0x8180a5CA4E3B94045e05A9313777955f7518D757"
	cat := "0x4A"
	l := NewCatLog(addr, cat)
	if err := self.storage.StoreCatLog(l); err != nil {
		return err
	}
	gotCat, err := self.storage.GetCategory(addr)
	if err != nil {
		return err
	}
	if gotCat != lowercaseCat {
		return errors.New(fmt.Sprintf("Got unexpected category. Expected(%s) Got(%s)",
			lowercaseCat, gotCat))
	}
	gotCat, err = self.storage.GetCategory(lowercaseAddr)
	if err != nil {
		return err
	}
	if gotCat != lowercaseCat {
		return errors.New(fmt.Sprintf("Got unexpected category. Expected(%s) Got(%s)",
			lowercaseCat, gotCat))
	}
	user, err := self.storage.GetUserOfAddress(lowercaseAddr)
	// initialy user is identical to the address
	if err != nil {
		return err
	}
	if user != lowercaseAddr {
		return errors.New(fmt.Sprintf("Got unexpected user. Expected(%s) Got(%s)",
			user, lowercaseAddr))
	}
	addresses, err := self.storage.GetAddressesOfUser(user)
	if err != nil {
		return err
	}
	if addresses[0] != lowercaseAddr {
		return errors.New(fmt.Sprintf("Got unexpected addresses. Expected(%v) Got(%v)",
			addresses, []string{lowercaseAddr}))
	}
	return nil
}

func (self *StorageTest) TestUpdateUserAddressesThenStoreCatLog() error {
	email := "victor@kyber.network"
	addr1 := "0x8180a5ca4e3b94045e05a9313777955f7518d757"
	addr2 := "0xcbac9e86e0b7160f1a8e4835ad01dd51c514afce"
	addr3 := "0x0ccd5bd8eb6822d357d7aef833274502e8b4b8ac"
	cat := "0x0000000000000000000000000000000000000000000000000000000000000004"

	self.storage.UpdateUserAddresses(
		email, []string{addr1, addr3},
	)
	// test if pending addresses are correct
	pendingAddrs, err := self.storage.GetPendingAddresses()
	if err != nil {
		return err
	}
	expectedAddresses := map[string]bool{
		addr1: true,
		addr3: true,
	}
	if len(pendingAddrs) != len(expectedAddresses) {
		return errors.New(
			fmt.Sprintf("Expected to get %d addresses, got %d addresses", len(expectedAddresses), len(pendingAddrs)))
	}
	for _, addr := range pendingAddrs {
		if _, found := expectedAddresses[addr]; !found {
			return errors.New(fmt.Sprintf("Expected to find %s, got not found", addr))
		}
	}
	self.storage.UpdateUserAddresses(
		email, []string{addr1, addr2},
	)
	// test if pending addresses are correct
	pendingAddrs, err = self.storage.GetPendingAddresses()
	if err != nil {
		return err
	}
	expectedAddresses = map[string]bool{
		addr1: true,
		addr2: true,
	}
	if len(pendingAddrs) != len(expectedAddresses) {
		return errors.New(
			fmt.Sprintf("Expected to get %d addresses, got %d addresses", len(expectedAddresses), len(pendingAddrs)))
	}
	for _, addr := range pendingAddrs {
		if _, found := expectedAddresses[addr]; !found {
			return errors.New(fmt.Sprintf("Expected to find %s, got not found", addr))
		}
	}
	// Start receiving cat logs
	self.storage.StoreCatLog(NewCatLog(addr1, cat))
	self.storage.UpdateUserAddresses(
		email, []string{addr1, addr2},
	)
	// test if pending addresses are correct
	pendingAddrs, err = self.storage.GetPendingAddresses()
	if err != nil {
		return err
	}
	expectedAddresses = map[string]bool{
		addr2: true,
	}
	if len(pendingAddrs) != len(expectedAddresses) {
		return errors.New(
			fmt.Sprintf("Expected to get %d addresses, got %d addresses", len(expectedAddresses), len(pendingAddrs)))
	}
	for _, addr := range pendingAddrs {
		if _, found := expectedAddresses[addr]; !found {
			return errors.New(fmt.Sprintf("Expected to find %s, got not found", addr))
		}
	}
	self.storage.StoreCatLog(NewCatLog(addr2, cat))

	gotAddresses, err := self.storage.GetAddressesOfUser(email)
	if err != nil {
		return err
	}
	// test addresses of user
	expectedAddresses = map[string]bool{
		addr1: true,
		addr2: true,
	}
	if len(gotAddresses) != len(expectedAddresses) {
		return errors.New(
			fmt.Sprintf("Expected to get %d addresses, got %d addresses", len(expectedAddresses), len(gotAddresses)))
	}
	for _, addr := range gotAddresses {
		if _, found := expectedAddresses[addr]; !found {
			return errors.New(fmt.Sprintf("Expected to find %s, got not found", addr))
		}
	}
	gotUser, err := self.storage.GetUserOfAddress(addr1)
	if err != nil {
		return err
	}
	if gotUser != email {
		return errors.New(fmt.Sprintf("Expected to get %s, got %s", email, gotUser))
	}
	gotUser, err = self.storage.GetUserOfAddress(addr2)
	if err != nil {
		return err
	}
	if gotUser != email {
		return errors.New(fmt.Sprintf("Expected to get %s, got %s", email, gotUser))
	}
	return nil
}

func (self *StorageTest) TestStoreCatLogThenUpdateUserAddresses() error {
	email := "Victor@kyber.network"
	lowercaseEmail := "victor@kyber.network"
	addr1 := "0x8180a5CA4E3B94045e05A9313777955f7518D757"
	lowercaseAddr1 := "0x8180a5ca4e3b94045e05a9313777955f7518d757"
	addr2 := "0xcbac9e86e0b7160f1a8e4835ad01dd51c514afce"
	cat := "0x4A"

	self.storage.StoreCatLog(NewCatLog(addr1, cat))
	self.storage.StoreCatLog(NewCatLog(addr2, cat))
	err := self.storage.UpdateUserAddresses(
		email, []string{addr1, addr2},
	)
	if err != nil {
		return err
	}
	gotAddresses, err := self.storage.GetAddressesOfUser(lowercaseEmail)
	if err != nil {
		return err
	}
	expectedAddresses := map[string]bool{
		lowercaseAddr1: true,
		addr2:          true,
	}
	if len(gotAddresses) != len(expectedAddresses) {
		return errors.New(
			fmt.Sprintf("Expected to get %d addresses, got %d addresses", len(expectedAddresses), len(gotAddresses)))
	}
	for _, addr := range gotAddresses {
		if _, found := expectedAddresses[addr]; !found {
			return errors.New(fmt.Sprintf("Expected to find %s, got not found", addr))
		}
	}
	gotUser, err := self.storage.GetUserOfAddress(addr1)
	if err != nil {
		return err
	}
	if gotUser != lowercaseEmail {
		return errors.New(fmt.Sprintf("Expected to get %s, got %s", lowercaseEmail, gotUser))
	}
	gotUser, err = self.storage.GetUserOfAddress(addr2)
	if err != nil {
		return err
	}
	if gotUser != lowercaseEmail {
		return errors.New(fmt.Sprintf("Expected to get %s, got %s", lowercaseEmail, gotUser))
	}
	return nil
}
