package exchange

import (
	"log"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/settings"
)

func getExchangePairsAndFeesFromConfig(
	addressConfig map[string]string,
	feeConfig common.ExchangeFees,
	minDepositConfig common.ExchangesMinDeposit,
	exchange string, sett *settings.Settings) ([]common.Token, []common.TokenPair, common.ExchangeFees, common.ExchangesMinDeposit) {

	tokens := []common.Token{}
	pairs := []common.TokenPair{}
	fees := common.ExchangeFees{
		feeConfig.Trading,
		common.FundingFee{
			map[string]float64{},
			map[string]float64{},
		},
	}
	minDeposit := common.ExchangesMinDeposit{}
	for tokenID := range addressConfig {
		token, err := sett.Tokens.GetInternalTokenByID(tokenID)
		if err != nil {
			log.Panicf("Must Get Internal Token failed :%s", err)
		}
		tokens = append(tokens, token)
		if tokenID != "ETH" {
			pair := sett.Tokens.MustCreateTokenPair(tokenID, "ETH")
			pairs = append(pairs, pair)
		}
		if _, exist := feeConfig.Funding.Withdraw[tokenID]; exist {
			fees.Funding.Withdraw[tokenID] = feeConfig.Funding.Withdraw[tokenID] * 2
		} else {
			panic(tokenID + " is not found in " + exchange + " withdraw fee config file")
		}
		if _, exist := feeConfig.Funding.Deposit[tokenID]; exist {
			fees.Funding.Deposit[tokenID] = feeConfig.Funding.Deposit[tokenID]
		} else {
			panic(tokenID + " is not found in " + exchange + " binance deposit fee config file")
		}
		log.Printf("minDepositConfig: %v", minDepositConfig)
		if _, exist := minDepositConfig[tokenID]; exist {
			minDeposit[tokenID] = minDepositConfig[tokenID] * 2
		} else {
			panic(tokenID + " is not found in " + exchange + " min deposit config file")
		}
	}
	return tokens, pairs, fees, minDeposit
}
