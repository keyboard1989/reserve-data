package exchange

type SimulatedLiquiEndpoint struct {
	*RealLiquiEndpoint
}

func NewSimulatedLiquiEndpoint(authendpoint string) *SimulatedLiquiEndpoint {
	realone := NewRealLiquiEndpoint()
	realone.AuthenticatedEndpoint = authendpoint
	return &SimulatedLiquiEndpoint{realone}
}
