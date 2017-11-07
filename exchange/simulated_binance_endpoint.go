package exchange

import (
	"os"
)

type SimulatedBinanceEndpoint struct {
	*RealBinanceEndpoint
}

func NewSimulatedBinanceEndpoint() *SimulatedBinanceEndpoint {
	realone := NewRealBinanceEndpoint()
	fakeAuthenticatedEndpoint := "http://127.0.0.1:8000/account"
	if len(os.Args) > 1 {
		fakeAuthenticatedEndpoint = os.Args[1]
	}
	realone.AuthenticatedEndpoint = fakeAuthenticatedEndpoint
	return &SimulatedBinanceEndpoint{realone}
}
