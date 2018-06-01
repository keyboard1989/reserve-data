package settings

import (
	"log"

	"github.com/KyberNetwork/reserve-data/common"
)

type Settings struct {
	Tokens *TokenSetting
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
		if common.RunningMode() == common.SIMULATION_MODE {
			tokenPath = simPath
		}
		if err = setting.LoadTokenFromFile(tokenPath); err != nil {
			log.Printf("Setting Init: Can not load Token from file: %s, Token DB is needed to be updated manually", err)
		}
	}
}
