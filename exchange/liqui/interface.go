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
	baseurl := "127.0.0.1"
	if len(os.Args) > 1 {
		baseurl = os.Args[1]
	}
	return baseurl + ":5000"
}

func (self *SimulatedInterface) PublicEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
	// return "https://api.liqui.io/api/3"
}

func (self *SimulatedInterface) AuthenticatedEndpoint(timepoint uint64) string {
	return fmt.Sprintf("%s?timestamp=%d", self.baseurl(), timepoint)
}

func NewSimulatedInterface() *SimulatedInterface {
	return &SimulatedInterface{}
}

type KovanInterface struct{}

func (self *KovanInterface) baseurl() string {
	baseurl := "http://simulator:8000"
	if len(os.Args) > 1 {
		baseurl = os.Args[1]
	}
	return baseurl
}

func (self *KovanInterface) PublicEndpoint(timepoint uint64) string {
	return "https://api.liqui.io/api/3"
}

func (self *KovanInterface) AuthenticatedEndpoint(timepoint uint64) string {
	return "https://api.liqui.io/tapi"
}

func NewKovanInterface() *KovanInterface {
	return &KovanInterface{}
}

type DevInterface struct{}

func (self *DevInterface) PublicEndpoint(timepoint uint64) string {
	// return "https://api.liqui.io/api/3"
	return "http://192.168.25.16:5000"
}

func (self *DevInterface) AuthenticatedEndpoint(timepoint uint64) string {
	// return "https://api.liqui.io/tapi"
	return "http://192.168.25.16:5000"
}

func NewDevInterface() *DevInterface {
	return &DevInterface{}
}
