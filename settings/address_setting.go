package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	settingstorage "github.com/KyberNetwork/reserve-data/settings/storage"
)

const (
	RESERVE         string = "reserve"
	BURNER          string = "burner"
	BANK            string = "bank"
	NETWORK         string = "network"
	WRAPPER         string = "wrapper"
	PRICING         string = "pricing"
	WHITELIST       string = "whitelist"
	SETRATE         string = "setrate"
	RESERVE3RDPARTY string = "third_party_reserves"
	OLDNETWORK      string = "oldNetworks"
	OLDBURNER       string = "oldBurners"
)

// AddressConfig type defines a list of address attribute avaiable in core.
// It is used mainly for
type AddressConfig struct {
	Bank               string   `json:"bank"`
	Reserve            string   `json:"reserve"`
	Network            string   `json:"network"`
	Wrapper            string   `json:"wrapper"`
	Pricing            string   `json:"pricing"`
	FeeBurner          string   `json:"feeburner"`
	Whitelist          string   `json:"whitelist"`
	ThirdPartyReserves []string `json:"third_party_reserves"`
	SetRate            string   `json:"setrate"`
}

// AddressSetting type defines component to handle all address setting in core.
// It contains the storage interface used to query addresses.
type AddressSetting struct {
	Storage AddressStorage
}

func NewAddressSetting(dbPath string) (*AddressSetting, error) {
	BoltAddressStorage, err := settingstorage.NewBoltAddressStorage(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Setting Init: Can not create bolt address storage (%s)", err)
	}
	addressSetting := AddressSetting{BoltAddressStorage}
	return &addressSetting, nil
}

func (setting *Settings) LoadAddressFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	addrs := AddressConfig{}
	if err = json.Unmarshal(data, &addrs); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(BANK, addrs.Bank); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(RESERVE, addrs.Reserve); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(NETWORK, addrs.Network); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(WRAPPER, addrs.Wrapper); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(PRICING, addrs.Pricing); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(BURNER, addrs.FeeBurner); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(WHITELIST, addrs.Whitelist); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress(SETRATE, addrs.SetRate); err != nil {
		return err
	}
	for _, addr := range addrs.ThirdPartyReserves {
		if err = setting.Address.Storage.AddAddressToSet(RESERVE3RDPARTY, addr); err != nil {
			return err
		}
	}
	return nil
}
