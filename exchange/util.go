package exchange

import (
	"log"

	"github.com/KyberNetwork/reserve-data/settings"

	"github.com/KyberNetwork/reserve-data/common"
)

func getExchangePairsAndFeesFromConfig(
	addressConfig map[string]string,
	minDepositConfig common.ExchangesMinDeposit,
	exchange settings.ExchangeName, setting Setting) ([]common.Token, []common.TokenPair, common.ExchangesMinDeposit, error) {

	tokens := []common.Token{}
	pairs := []common.TokenPair{}
	minDeposit := common.ExchangesMinDeposit{}

	for tokenID := range addressConfig {
		token, err := setting.GetInternalTokenByID(tokenID)
		if err != nil {
			log.Panicf("Must Get Internal Token failed :%s", err)
		}
		tokens = append(tokens, token)
		if tokenID != "ETH" {
			pair := setting.MustCreateTokenPair(tokenID, "ETH")
			pairs = append(pairs, pair)
		}
		log.Printf("minDepositConfig: %v", minDepositConfig)
		if _, exist := minDepositConfig[tokenID]; exist {
			minDeposit[tokenID] = minDepositConfig[tokenID] * 2
		} else {
			panic(tokenID + " is not found in " + string(exchange) + " min deposit config file")
		}
	}
	return tokens, pairs, minDeposit, nil
}
