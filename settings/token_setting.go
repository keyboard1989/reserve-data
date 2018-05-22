package settings

import (
	"encoding/json"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-data/common"
)

type TokenConfig struct {
	Tokens map[string]common.Token `json:"tokens"`
}

type TokenSetting struct {
	Storage TokenStorage
}

func NewTokenSetting(storage TokenStorage) *TokenSetting {
	return &TokenSetting{storage}
}

func UpdateToken(t common.Token) error {
	if err := setting.Tokens.Storage.AddTokenByID(t); err != nil {
		return err
	}
	if err := setting.Tokens.Storage.AddTokenByAddress(t); err != nil {
		return err
	}
	return nil
}

func LoadTokenFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	tokens := TokenConfig{}
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return err
	}
	for _, t := range tokens.Tokens {
		err = UpdateToken(t)
		if err != nil {
			return err
		}
	}
	return nil
}
