package exchange

type SimulatedLiquiEndpoint struct {
	*RealLiquiEndpoint
}

func NewSimulatedLiquiEndpoint() *SimulatedLiquiEndpoint {
	realone := NewRealLiquiEndpoint()
	realone.AuthenticatedEndpoint = "http://127.0.0.1:8000"
	return &SimulatedLiquiEndpoint{realone}
}
