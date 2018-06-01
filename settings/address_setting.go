package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"

	settingstorage "github.com/KyberNetwork/reserve-data/settings/storage"
)

const (
	addressDBFileName            string = "address.db"
	ADDRES_DEFAULT_JSON_PATH     string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json"
	ADDRES_DEFAULT_JSON_SIM_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json"
)

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

type AddressSetting struct {
	Storage AddressStorage
}

func NewAddressSetting(dbPath string) *AddressSetting {
	BoltAddressStorage, err := settingstorage.NewBoltAddressStorage(dbPath)
	if err != nil {
		log.Panicf("Setting Init: Can not create bolt address storage (%s)", err)
	}
	addressSetting := AddressSetting{BoltAddressStorage}
	return &addressSetting
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
	if err = setting.Address.Storage.UpdateOneAddress("bank", addrs.Bank); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("reserve", addrs.Reserve); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("network", addrs.Network); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("wrapper", addrs.Wrapper); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("pricing", addrs.Pricing); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("burner", addrs.FeeBurner); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("whitelist", addrs.Whitelist); err != nil {
		return err
	}
	if err = setting.Address.Storage.UpdateOneAddress("setrate", addrs.SetRate); err != nil {
		return err
	}
	for _, addr := range addrs.ThirdPartyReserves {
		if err = setting.Address.Storage.AddAddressToSet("third_party_reserves", addr); err != nil {
			return err
		}
	}
	return nil
}
