package exchange

import (
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
)

func getExchangePairsAndFeesFromConfig(
	exchange settings.ExchangeName, setting Setting) ([]common.Token, error) {

	tokens := []common.Token{}
	addressConfig, err := setting.GetDepositAddress(exchange)
	if err != nil {
		return nil, fmt.Errorf("get exchange Deposit Address failed: (%s)", err)
	}
	for tokenID := range addressConfig {
		token, err := setting.GetInternalTokenByID(tokenID)
		if err != nil {
			return nil, fmt.Errorf("Must Get Internal Token failed :%s", err)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
