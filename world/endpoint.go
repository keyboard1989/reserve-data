package world

type Endpoint interface {
	GoldDataEndpoint() string
}

type RealEndpoint struct {
}

func (self RealEndpoint) GoldDataEndpoint() string {
	return "https://datafeed.digix.global/tick/"
}

type SimulatedEndpoint struct {
}

func (self SimulatedEndpoint) GoldDataEndpoint() string {
	return ""
}
