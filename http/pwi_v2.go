package http

import (
	"encoding/json"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

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
		println(len(et))
		return false
	}

	for _, field := range requiredFields {
		println(field)
		if _, ok := et[field]; !ok {
			return false
		}
	}
	return true
}

// PWIEquationV2Input is the input SetPWIEquationV2 api.
type PWIEquationV2Input map[string]PWIEquationTokenV2

// isValid validates the input instance and return true if it is valid.
// Example input:
// [{"token_id": {equation_token}}, ...]
func (input PWIEquationV2Input) isValid() bool {
	for _, et := range input {
		if !et.isValid() {
			return false
		}
	}
	return true
}

// SetPWIEquationV2 stores the given PWI equations to pending for later evaluation.
func (self *HTTPServer) SetPWIEquationV2(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{}, []Permission{ConfigurePermission})
	if !ok {
		return
	}

	sampleInput := &PWIEquationV2Input{
		"EOS": PWIEquationTokenV2{
			"bid": PWIEquationV2{
				A: 1,
				B: 2,
				C: 3,
			},
		},
	}

	encoded, _ := json.Marshal(sampleInput)
	println(string(encoded))

	data := []byte(postForm.Get("data"))
	if len(data) > MAX_DATA_SIZE {
		httputil.ResponseFailure(c, httputil.WithError(errDataSizeExceed))
		return
	}

	var input PWIEquationV2Input
	if err := json.Unmarshal(data, &input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	if !input.isValid() {
		httputil.ResponseFailure(c, httputil.WithReason("invalid input"))
		return
	}

	if err := self.metric.StorePendingPWIEquationV2(data); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}
