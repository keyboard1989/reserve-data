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
	baseurl := "http://127.0.0.1"
	if len(os.Args) > 1 {
		baseurl = os.Args[1]
	}
	return baseurl + ":5300"
}

func (self *SimulatedInterface) PublicEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s/api/%s/public?timestamp=%d", self.baseurl(), apiVersion, timepoint)
}

func (self *SimulatedInterface) MarketEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s/api/%s/market?timestamp=%d", self.baseurl(), apiVersion, timepoint)
}

func (self *SimulatedInterface) AccountEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s/api/%s/account?timestamp=%d", self.baseurl(), apiVersion, timepoint)
}

func NewSimulatedInterface() *SimulatedInterface {
	return &SimulatedInterface{}
}

type DevInterface struct{}

func (self *DevInterface) PublicEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://bittrex.com/api/" + apiVersion + "/public"
}

func (self *DevInterface) MarketEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://bittrex.com/api/" + apiVersion + "/market"
}

func (self *DevInterface) AccountEndpoint(timepoint uint64) string {
	// ignore timepoint because timepoint is only relevant in simulation
	return "https://bittrex.com/api/" + apiVersion + "/account"
}

func NewDevInterface() *DevInterface {
	return &DevInterface{}
}
