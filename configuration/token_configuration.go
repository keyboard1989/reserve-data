package configuration

import "github.com/KyberNetwork/reserve-data/common"

type TokenConfiguration struct {
	storage TokenStorage
}

func NewTokenConfiguration(storage TokenStorage) *TokenConfiguration {
	return &TokenConfiguration{storage}
}

func (self *TokenConfiguration) AddToken(t common.Token, active bool, knSupported bool) error {
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
