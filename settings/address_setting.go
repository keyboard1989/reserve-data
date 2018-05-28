package settings

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type AddressConfig struct {
	Exchanges          map[string]common.Exchange `json:"exchanges"`
	Bank               string                     `json:"bank"`
	Reserve            string                     `json:"reserve"`
	Network            string                     `json:"network"`
	Wrapper            string                     `json:"wrapper"`
	Pricing            string                     `json:"pricing"`
	FeeBurner          string                     `json:"feeburner"`
	Whitelist          string                     `json:"whitelist"`
	ThirdPartyReserves []string                   `json:"third_party_reserves"`
	Intermediator      string                     `json:"intermediator"`
	SetRate            string                     `json:"setrate"`
}
