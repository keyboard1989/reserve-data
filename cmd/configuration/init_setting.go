package configuration

import (
	"path/filepath"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
)

const (
	tokenDBFileName           string = "token.db"
	tokenSettingDefaultName   string = "mainnet_setting.json"
	tokenSettingSimName       string = "shared/deployment_dev.json"
	addressDBFileName         string = "address.db"
	addressSettingDefaultName string = "mainnet_setting.json"
	addressSettingSimName     string = "deployment_dev.json"
)

func createSetting() *settings.Settings {
	tokensSetting := settings.NewTokenSetting(filepath.Join(common.CmdDirLocation(), tokenDBFileName))
	addressSetting := settings.NewAddressSetting(filepath.Join(common.CmdDirLocation(), addressDBFileName))
	setting := settings.NewSetting(tokensSetting, addressSetting)
	setting.HandleEmptyToken(filepath.Join(common.CmdDirLocation(), tokenSettingDefaultName), filepath.Join(common.CmdDirLocation(), tokenSettingSimName))
	setting.HandleEmptyAddress(filepath.Join(common.CmdDirLocation(), addressSettingDefaultName), filepath.Join(common.CmdDirLocation(), addressSettingSimName))
	return setting
}
