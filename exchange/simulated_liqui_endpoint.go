package exchange

import (
	"os"
)

type SimulatedLiquiEndpoint struct {
	*RealLiquiEndpoint
}

func NewSimulatedLiquiEndpoint() *SimulatedLiquiEndpoint {
	realone := NewRealLiquiEndpoint()
	fakeEndpoint := "http://127.0.0.1:8000"
	if len(os.Args) > 1 {
		fakeEndpoint = os.Args[1]
	}
	realone.AuthenticatedEndpoint = fakeEndpoint
	return &SimulatedLiquiEndpoint{realone}
}
