package signer

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"

	ethereum "github.com/ethereum/go-ethereum/common"
)

type FileSigner struct {
	LiquiKey    string `json:"liqui_key"`
	LiquiSecret string `json:"liqui_secret"`
}

func (self FileSigner) GetLiquiKey() string {
	return self.LiquiKey
}

func (self FileSigner) LiquiSign(msg string) string {
	mac := hmac.New(sha512.New, []byte(self.LiquiSecret))
	mac.Write([]byte(msg))
	return ethereum.Bytes2Hex(mac.Sum(nil))
}

func NewFileSigner(file string) *FileSigner {
	raw, err := ioutil.ReadFile(file)
	fmt.Printf("read from file: %s\n", raw)
	if err != nil {
		panic(err)
	}
	signer := FileSigner{}
	err = json.Unmarshal(raw, &signer)
	if err != nil {
		panic(err)
	}
	return &signer
}
