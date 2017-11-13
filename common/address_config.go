package common

import (
	"encoding/json"
	"io/ioutil"
)

type token struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Decimals int64  `json:"decimals"`
}

type AddressConfig struct {
	Tokens    map[string]token  `json:"tokens"`
	Exchanges map[string]string `json:"exchanges"`
	Bank      string            `json:"bank"`
	Reserve   string            `json:"reserve"`
	Network   string            `json:"network"`
	Wrapper   string            `json:"wrapper"`
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
