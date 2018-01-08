package okex

import (
	"os"
)

type Interface interface {
	PublicEnpoint() string
	AuthenticatedEndpoint() string
}

type RealInterface struct{}

func (self *RealInterface) PublicEnpoint() string {
	return "https://www.okex.com/api/v1/"
}

func (self *RealInterface) AuthenticatedEndpoint() string {
	return "https://www.okex.com/api/v1/"
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

func (self *SimulatedInterface) PublicEnpoint() string {
	return "https://www.okex.com/api/v1/"
}

func (self *SimulatedInterface) AuthenticatedEndpoint() string {
	return self.baseurl()
}

func NewSimulatedInterface() *SimulatedInterface {
	return &SimulatedInterface{}
}

type KovanInterface struct{}

func (self *KovanInterface) baseurl() string {
	baseurl := "127.0.0.1"
	if len(os.Args) > 1 {
		baseurl = os.Args[1]
	}
	return baseurl + ":5100"
}

func (self *KovanInterface) PublicEnpoint() string {
	return "https://www.okex.com/api/v1/"
}

func (self *KovanInterface) AuthenticatedEndpoint() string {
	return self.baseurl()
}

func NewKovanInterface() *KovanInterface {
	return &KovanInterface{}
}

type DevInterface struct{}

func (self *DevInterface) PublicEnpoint() string {
	return "https://www.okex.com/api/v1/"
}

func (self *DevInterface) AuthenticatedEndpoint() string {
	return "https://www.okex.com/api/v1/"
}

func NewDevInterface() *DevInterface {
	return &DevInterface{}
}
