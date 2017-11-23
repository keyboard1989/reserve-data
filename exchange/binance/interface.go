package binance

import (
	"os"
)

type Interface interface {
	PublicEndpoint() string
	AuthenticatedEndpoint() string
}

type RealInterface struct{}

func (self *RealInterface) PublicEndpoint() string {
	return "https://www.binance.com"
}

func (self *RealInterface) AuthenticatedEndpoint() string {
	return "https://www.binance.com"
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
	return baseurl + ":5100"
}

func (self *SimulatedInterface) PublicEndpoint() string {
	return self.baseurl()
}

func (self *SimulatedInterface) AuthenticatedEndpoint() string {
	return self.baseurl()
}

func NewSimulatedInterface() *SimulatedInterface {
	return &SimulatedInterface{}
}

type KovanInterface struct{}

func (self *KovanInterface) PublicEndpoint() string {
	// return "http://192.168.24.197:5100"
	return "https://www.binance.com"
}

func (self *KovanInterface) AuthenticatedEndpoint() string {
	// return "http://192.168.24.197:5100"
	return "https://www.binance.com"
}

func NewKovanInterface() *KovanInterface {
	return &KovanInterface{}
}

type DevInterface struct{}

func (self *DevInterface) PublicEndpoint() string {
	return "https://www.binance.com"
}

func (self *DevInterface) AuthenticatedEndpoint() string {
	return "https://www.binance.com"
}

func NewDevInterface() *DevInterface {
	return &DevInterface{}
}
