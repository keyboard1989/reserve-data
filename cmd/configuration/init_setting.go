package configuration

import (
	"path/filepath"

	"github.com/KyberNetwork/reserve-data/cmd/configuration/mode"
	"github.com/KyberNetwork/reserve-data/settings"
)

const (
	tokenDBFileName             string = "token.db"
	TOKEN_DEFAULT_JSON_PATH     string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/mainnet_setting.json"
	TOKEN_DEFAULT_JSON_SIM_PATH string = "/go/src/github.com/KyberNetwork/reserve-data/cmd/shared/deployment_dev.json"
)

func createSetting() *settings.Settings {
	tokensSetting := settings.NewTokenSetting(filepath.Join(mode.CmdDirLocation(), tokenDBFileName))
	setting := settings.NewSetting(tokensSetting)
	setting.HandleEmptyToken(TOKEN_DEFAULT_JSON_PATH, TOKEN_DEFAULT_JSON_SIM_PATH)
	return setting
}
