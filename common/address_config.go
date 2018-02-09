package common

import (
	"encoding/json"
	"io/ioutil"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type token struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Decimals int64  `json:"decimals"`
}

type exchange map[string]string

type AddressConfig struct {
	Tokens    map[string]token    `json:"tokens"`
	Exchanges map[string]exchange `json:"exchanges"`
	Bank      string              `json:"bank"`
	Reserve   string              `json:"reserve"`
	Network   string              `json:"network"`
	Wrapper   string              `json:"wrapper"`
	Pricing   string              `json:"pricing"`
	FeeBurner string              `json:"feeburner"`
}

func GetAddressConfigFromFile(path string) (AddressConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return AddressConfig{}, err
	} else {
		result := AddressConfig{}
		err := json.Unmarshal(data, &result)
		return result, err
	}
}

type Addresses struct {
	Tokens      map[string]ethereum.Address   `json:"tokens"`
	Exchanges   map[ExchangeID]TokenAddresses `json:"exchanges"`
	WrapperAddress   ethereum.Address         `json:"wrapper"`
	PricingAddress   ethereum.Address         `json:"pricing"`
	ReserveAddress   ethereum.Address         `json:"reserve"`
	FeeBurnerAddress ethereum.Address         `json:"feeburner"`
	NetworkAddress   ethereum.Address         `json:"network"`
}

type TokenAddresses map[string]ethereum.Address
