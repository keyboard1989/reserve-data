package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/KyberNetwork/reserve-data/cmd/configuration/mode"
	"github.com/KyberNetwork/reserve-data/common"
	settingstorage "github.com/KyberNetwork/reserve-data/settings/storage"
)

const (
	TOKEN_DB_FILE_PATH          string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/token.db"
	TOKEN_DEFAULT_JSON_PATH     string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json"
	TOKEN_DEFAULT_JSON_SIM_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json"
)

type token struct {
	Address                 string `json:"address"`
	Name                    string `json:"name"`
	Decimal                 int64  `json:"decimals"`
	Active                  bool   `json:"internal use"`
	Internal                bool   `json:"listed"`
	MinimalRecordResolution string `json:"minimalRecordResolution"`
	MaxTotalImbalance       string `json:"maxPerBlockImbalance"`
	MaxPerBlockImbalance    string `json:"maxTotalImbalance"`
}

type TokenConfig struct {
	Tokens map[string]token `json:"tokens"`
}

type TokenSetting struct {
	Storage TokenStorage
}

func NewTokenSetting() *TokenSetting {
	BoltTokenStorage, err := settingstorage.NewBoltTokenStorage(filepath.Join(mode.CmdDirLocation(), tokenDBFileName))
	if err != nil {
		log.Panicf("Setting Init: Can not create bolt token storage (%s)", err)
	}
	tokenSetting := TokenSetting{BoltTokenStorage}
	return &tokenSetting

}

func UpdateToken(t common.Token) error {
	return setting.Tokens.Storage.UpdateToken(t)
}

func LoadTokenFromFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	tokens := TokenConfig{}
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &tokens); err != nil {
		return err
	}
	for id, t := range tokens.Tokens {
		token := common.NewToken(id, t.Name, t.Address, t.Decimal, t.Active, t.Internal, t.MinimalRecordResolution, t.MaxPerBlockImbalance, t.MaxTotalImbalance)
		if err = UpdateToken(token); err != nil {
			return err
		}
	}
	return nil
}
