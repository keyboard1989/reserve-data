package settings

import (
	"encoding/json"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-data/common"
)

type token struct {
	Address          string `json:"address"`
	Name             string `json:"name"`
	Decimals         int64  `json:"decimals"`
	KNReserveSupport bool   `json:"internal use"`
	Active           bool   `json:"listed"`
}

type TokenConfig struct {
	Tokens map[string]token `json:"tokens"`
}

type TokenSetting struct {
	storage TokenStorage
}

func NewTokenSetting(storage TokenStorage) *TokenSetting {
	return &TokenSetting{storage}
}

func (self *TokenSetting) AddToken(t common.Token, active bool, knSupported bool) error {
	if err := self.storage.AddTokenByID(t); err != nil {
		return err
	}
	if err := self.storage.AddTokenByAddress(t); err != nil {
		return err
	}
	if active {
		if err := self.storage.AddActiveTokenByID(t); err != nil {
			return err
		}
		if err := self.storage.AddActiveTokenByAddress(t); err != nil {
			return err
		}
		if knSupported {
			if err := self.storage.AddInternalActiveTokenByID(t); err != nil {
				return err
			}
			if err := self.storage.AddInternalActiveTokenByAddress(t); err != nil {
				return err
			}
		} else {
			if err := self.storage.AddExternalActiveTokenByID(t); err != nil {
				return err
			}
			if err := self.storage.AddExternalActiveTokenByAddress(t); err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *TokenSetting) LoadTokenFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	tokens := TokenConfig{}
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return err
	}
	for id, t := range tokens.Tokens {
		tok := common.Token{id, t.Address, t.Decimals}
		err = self.AddToken(tok, t.Active, t.KNReserveSupport)
		if err != nil {
			return err
		}
	}
	return nil
}
