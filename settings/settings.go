package settings

import (
	"log"
	"os"
)

type Settings struct {
	Tokens  *TokenSetting
	Address *AddressSetting
}

func NewSetting(tokenSetting *TokenSetting) *Settings {
	setting := Settings{tokenSetting}
	return &setting
}

func (setting *Settings) HandleEmptyToken(normalPath, simPath string) {
	allToks, err := setting.GetAllTokens()
	if err != nil || len(allToks) < 1 {
		if err != nil {
			log.Printf("Setting Init: Token DB is faulty (%s), attempt to load token from file", err)
		} else {
			log.Printf("Setting Init: Token DB is empty, attempt to load token from file")
		}
		tokenPath := normalPath
		if os.Getenv("KYBER_ENV") == "simulation" {
			tokenPath = simPath
		}

		if err = setting.LoadTokenFromFile(tokenPath); err != nil {
			log.Printf("Setting Init: Can not load Token from file: %s, Token DB is needed to be updated manually", err)
		}
	}
}

func handleEmptyAddress() {
	addressCount, err := setting.Address.Storage.CountAddress()
	if addressCount == 0 || err != nil {
		if err != nil {
			log.Printf("Setting Init: Address DB is faulty (%s), attempt to load Address from file", err)
		} else {
			log.Printf("Setting Init: Address DB is empty, attempt to load address from file")
		}
		addressPath := ADDRES_DEFAULT_JSON_PATH
		if os.Getenv("KYBER_ENV") == "simulation" {
			addressPath = ADDRES_DEFAULT_JSON_SIM_PATH
		}
		if err = LoadAddressFromFile(addressPath); err != nil {
			log.Printf("Setting Init: Can not load Address from file: %s, address DB is needed to be updated manually", err)
		}
	}
}
