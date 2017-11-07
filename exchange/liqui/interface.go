package liqui

import (
	"fmt"
	"os"
)

type Interface interface {
	PublicEndpoint(timepoint uint64) string
	AuthenticatedEndpoint(timepoint uint64) string
}

type RealInterface struct{}

func (self *RealInterface) PublicEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://api.liqui.io/api/3"
}

func (self *RealInterface) AuthenticatedEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://api.liqui.io/tapi"
}

func NewRealInterface() *RealInterface {
	return &RealInterface{}
}

type SimulatedInterface struct{}

func (self *SimulatedInterface) baseurl() string {
	baseurl := "http://127.0.0.1:8000"
	if len(os.Args) > 1 {
		baseurl = os.Args[1]
	}
	return baseurl
}

func (self *SimulatedInterface) PublicEndpoint(timepoint uint64) string {
	// return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
	return "https://api.liqui.io/api/3"
}

func (self *SimulatedInterface) AuthenticatedEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
}

func NewSimulatedInterface() *SimulatedInterface {
	return &SimulatedInterface{}
}
