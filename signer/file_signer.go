package signer

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type FileSigner struct {
	LiquiKey    string `json:"liqui_key"`
	LiquiSecret string `json:"liqui_secret"`
	Keystore    string `json:"keystore_path"`
	Passphrase  string `json:"passphrase"`
	opts        *bind.TransactOpts
}

func (self FileSigner) GetAddress() ethereum.Address {
	return self.opts.From
}

func (self FileSigner) Sign(address ethereum.Address, tx *types.Transaction) (*types.Transaction, error) {
	return self.opts.Signer(types.HomesteadSigner{}, address, tx)
}

func (self FileSigner) GetTransactOpts() *bind.TransactOpts {
	return self.opts
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
	fmt.Printf("keystore: %s\n", signer.Keystore)
	keyio, err := os.Open(signer.Keystore)
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewTransactor(keyio, signer.Passphrase)
	if err != nil {
		panic(err)
	}

	auth.GasLimit = big.NewInt(1000000)
	auth.GasPrice = big.NewInt(20000000000)
	signer.opts = auth
	return &signer
}
