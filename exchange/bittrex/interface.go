package bittrex

import (
	"fmt"
	"os"
)

type Interface interface {
	PublicEndpoint(timepoint uint64) string
	MarketEndpoint(timepoint uint64) string
	AccountEndpoint(timepoint uint64) string
}

type RealInterface struct{}

const apiVersion string = "v1.1"

func (self *RealInterface) PublicEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://bittrex.com/api/" + apiVersion + "/public"
}

func (self *RealInterface) MarketEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://bittrex.com/api/" + apiVersion + "/market"
}

func (self *RealInterface) AccountEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://bittrex.com/api/" + apiVersion + "/account"
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
	return "https://bittrex.com/api/" + apiVersion + "/public"
	return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
}

func (self *SimulatedInterface) MarketEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
}

func (self *SimulatedInterface) AccountEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
}

func NewSimulatedInterface() *SimulatedInterface {
	return &SimulatedInterface{}
}
