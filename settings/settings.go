package settings

import (
	"log"

	settingstorage "github.com/KyberNetwork/reserve-data/settings/storage"
)

type Settings struct {
	Tokens  *TokenSetting
	Address *AddressSetting
}

// HandleEmptyToken will load the token settings from default file if the
// database is empty.
func WithHandleEmptyToken(pathJSON string) SettingOption {
	return func(setting *Settings) {
		allToks, err := setting.GetAllTokens()
		if err != nil || len(allToks) < 1 {
			if err != nil {
				log.Printf("Setting Init: Token DB is faulty (%s), attempt to load token from file", err)
			} else {
				log.Printf("Setting Init: Token DB is empty, attempt to load token from file")
			}
			if err = setting.LoadTokenFromFile(pathJSON); err != nil {
				log.Printf("Setting Init: Can not load Token from file: %s, Token DB is needed to be updated manually", err)
			}
		}
	}
}

// HandleEmptyAddress will load the address settings from default file if the
// database is empty.
func WithHandleEmptyAddress(pathJSON string) SettingOption {
	return func(setting *Settings) {
		addressCount, err := setting.Address.Storage.CountAddress()
		if addressCount == 0 || err != nil {
			if err != nil {
				log.Printf("Setting Init: Address DB is faulty (%s), attempt to load Address from file", err)
			} else {
				log.Printf("Setting Init: Address DB is empty, attempt to load address from file")
			}
			if err = setting.LoadAddressFromFile(pathJSON); err != nil {
				log.Printf("Setting Init: Can not load Address from file: %s, address DB is needed to be updated manually", err)
			}
		}
	}
}

// SettingOption sets the initialization behavior of the Settings instance.
type SettingOption func(s *Settings)

func NewSetting(dbPath string, options ...SettingOption) (*Settings, error) {
	boltSettingStorage, err := settingstorage.NewBoltSettingStorage(dbPath)
	if err != nil {
		return nil, err
	}
	tokenSetting, err := NewTokenSetting(boltSettingStorage)
	if err != nil {
		return nil, err
	}
	addressSetting, err := NewAddressSetting(boltSettingStorage)
	if err != nil {
		return nil, err
	}
	setting := &Settings{Tokens: tokenSetting,
		Address: addressSetting}
	for _, option := range options {
		option(setting)
	}
	return nil, setting
}
