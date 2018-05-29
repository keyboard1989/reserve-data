package metric

import (
	"log"

	"github.com/KyberNetwork/reserve-data/settings"
)

type TokenMetric struct {
	AfpMid float64
	Spread float64
}

type MetricEntry struct {
	Timestamp uint64
	// data contain all token metric for all tokens
	Data map[string]TokenMetric
}

type TokenMetricResponse struct {
	Timestamp uint64
	AfpMid    float64
	Spread    float64
}

// Metric list for one token
type MetricList []TokenMetricResponse

type MetricResponse struct {
	Timestamp  uint64
	ReturnTime uint64
	Data       map[string]MetricList
}

type TokenTargetQty struct {
	ID        uint64
	Timestamp uint64
	Data      string
	Status    string
	Type      int64
}

type PWIEquation struct {
	ID   uint64 `json:"id"`
	Data string `json:"data"`
}

type RebalanceControl struct {
	Status bool `json:"status"`
}

type SetrateControl struct {
	Status bool `json:"status"`
}

type TargetQtySet struct {
	TotalTarget        float64
	ReserveTarget      float64
	RebalanceThreshold float64
	TransferThreshold  float64
}

type TargetQtyStruct struct {
	SetTarget TargetQtySet
}

// PWIEquationV2 contains the information of a PWI equation.
type PWIEquationV2 struct {
	A                   float64
	B                   float64
	C                   float64
	MinMinSpread        float64
	PriceMultiplyFactor float64
}

// PWIEquationTokenV2 is a mapping between a token id and a PWI equation.
type PWIEquationTokenV2 map[string]PWIEquationV2

// isValid validates the input instance and return true if it is valid.
// Example:
// {
//  "bid": {
//    "a": "750",
//    "b": "500",
//    "c": "0",
//    "min_min_spread": "0",
//    "price_multiply_factor": "0"
//  },
//  "ask": {
//    "a": "800",
//    "b": "600",
//    "c": "0",
//    "min_min_spread": "0",
//    "price_multiply_factor": "0"
//  }
//}
func (et PWIEquationTokenV2) isValid() bool {
	var requiredFields = []string{"bid", "ask"}
	// validates that both bid and ask are present
	if len(et) != len(requiredFields) {
		return false
	}

	for _, field := range requiredFields {
		if _, ok := et[field]; !ok {
			return false
		}
	}
	return true
}

// PWIEquationRequestV2 is the input SetPWIEquationV2 api.
type PWIEquationRequestV2 map[string]PWIEquationTokenV2

// IsValid validates the input instance and return true if it is valid.
// Example input:
// [{"token_id": {equation_token}}, ...]
func (input PWIEquationRequestV2) IsValid() bool {
	for tokenID, et := range input {
		if !et.isValid() {
			return false
		}

		if _, err := settings.GetInternalTokenByID(tokenID); err != nil {
			log.Printf("unsupported token %s", tokenID)
			return false
		}
	}

	return true
}
