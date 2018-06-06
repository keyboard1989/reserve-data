package exchange

import (
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
)

func getExchangePairsAndFeesFromConfig(
	exchange settings.ExchangeName, setting Setting) ([]common.Token, []common.TokenPair, error) {

	tokens := []common.Token{}
	pairs := []common.TokenPair{}
	addressConfig, err := setting.GetDepositAddress(exchange)
	if err != nil {
		return nil, nil, fmt.Errorf("get exchange Deposit Address failed: (%s)", err)
	}
	for tokenID := range addressConfig {
		token, err := setting.GetInternalTokenByID(tokenID)
		if err != nil {
			return nil, nil, fmt.Errorf("Must Get Internal Token failed :%s", err)
		}
		tokens = append(tokens, token)
		if tokenID != "ETH" {
			pair := setting.MustCreateTokenPair(tokenID, "ETH")
			pairs = append(pairs, pair)
		}
	}
	return tokens, pairs, nil
}
