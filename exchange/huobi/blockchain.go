package huobi

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/rpc"
	"os"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type tbindex struct {
	BulkIndex   uint64
	IndexInBulk uint64
}

type NonceCorpus interface {
	GetAddress() ethereum.Address
	GetNextNonce() (*big.Int, error)
	MinedNonce() (*big.Int, error)
}

type Signer interface {
	GetAddress() ethereum.Address
	Sign(*types.Transaction) (*types.Transaction, error)
	GetTransactOpts() *bind.TransactOpts
}

type Blockchain struct {
	rpcClient          *rpc.Client
	client             *ethclient.Client
	intermediateSigner Signer
	tokens             []common.Token
	tokenIndices       map[string]tbindex
	nonceIntermediate  NonceCorpus
	chainType          string
}

func getNextNonce(n NonceCorpus) (*big.Int, error) {
	var nonce *big.Int
	var err error
	for i := 0; i < 3; i++ {
		nonce, err = n.GetNextNonce()
		if err == nil {
			return nonce, nil
		}
	}
	return nonce, err
}

func donothing() {}

func (self *Blockchain) getIntermediateTransactOpts(nonce, gasPrice *big.Int) (*bind.TransactOpts, context.CancelFunc, error) {
	shared := self.intermediateSigner.GetTransactOpts()
	var err error
	if nonce == nil {
		nonce, err = getNextNonce(self.nonceIntermediate)
	}
	if err != nil {
		return nil, donothing, err
	}
	if gasPrice == nil {
		gasPrice = big.NewInt(50100000000)
	}
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	result := bind.TransactOpts{
		shared.From,
		nonce,
		shared.Signer,
		shared.Value,
		gasPrice,
		shared.GasLimit,
		timeout,
	}
	return &result, cancel, nil
}

func packData(method string, params ...interface{}) ([]byte, error) {
	file, err := os.Open(
		"/go/src/github.com/KyberNetwork/reserve-data/blockchain/ERC20.abi")
	if err != nil {
		return nil, err
	}
	packabi, err := abi.JSON(file)
	if err != nil {
		return nil, err
	}
	data, err := packabi.Pack(method, params...)
	if err != nil {
		log.Println("Intermediator: Can not pack the data")
		return nil, err
	}
	return data, nil
}

func (self *Blockchain) SendTokenFromAccountToExchange(amount *big.Int, exchangeAddress ethereum.Address, tokenAddress ethereum.Address) (*types.Transaction, error) {
	opts, cancel, err := self.getIntermediateTransactOpts(nil, nil)
	ctx := opts.Context
	defer cancel()
	data, err := packData("transfer", exchangeAddress, amount)
	if err != nil {
		return nil, err
	}

	msg := ether.CallMsg{From: opts.From, To: &tokenAddress, Value: big.NewInt(0), Data: data}
	log.Printf("msg in hex is  %x", msg)
	gasLimit, err := self.client.EstimateGas(ensureContext(opts.Context), msg)
	if err != nil {
		log.Printf("Intermediator: Can not estimate gas %v", err)
		gasLimit = big.NewInt(25000)
	} else {
		log.Println("Intermediator: gas limit estimated is : %d", gasLimit)
		if (gasLimit.Cmp(big.NewInt(25000))) < 0 {
			gasLimit = big.NewInt(25000)
		}
	}
	//build tx, sign and send

	tx := types.NewTransaction(opts.Nonce.Uint64(), tokenAddress, big.NewInt(0), gasLimit, opts.GasPrice, data)
	signTX, err := self.intermediateSigner.Sign(tx)
	if err != nil {
		log.Println("Intermediator: Can not sign the transaction")
		return nil, err
	}
	err = self.client.SendTransaction(ctx, signTX)
	if err != nil {
		log.Println("Intermediator: ERROR: Can't send the transaction")
		return nil, err
	}
	return signTX, nil
}

func (self *Blockchain) SendETHFromAccountToExchange(amount *big.Int, exchangeAddress ethereum.Address) (*types.Transaction, error) {

	opts, cancel, err := self.getIntermediateTransactOpts(nil, nil)
	ctx := opts.Context
	defer cancel()
	//build msg and get gas limit
	msg := ether.CallMsg{From: opts.From, To: &exchangeAddress, Value: amount, Data: nil}
	gasLimit, err := self.client.EstimateGas(ensureContext(opts.Context), msg)
	//build tx, sign and send
	tx := types.NewTransaction(opts.Nonce.Uint64(), exchangeAddress, amount, gasLimit, opts.GasPrice, nil)
	signTX, err := self.intermediateSigner.Sign(tx)
	if err != nil {
		log.Println("Intermediator: Can not sign the transaction")
		return nil, err
	}
	log.Printf("Intermediator: the signed TX is: \n%v ", signTX)
	err = self.client.SendTransaction(ctx, signTX)
	if err != nil {
		log.Println("ERROR: Can't send the transaction")
		return nil, err
	}
	return signTX, nil
}

func (self *Blockchain) TransactionByHash(ctx context.Context, hash ethereum.Hash) (tx *rpcTransaction, isPending bool, err error) {
	var json *rpcTransaction
	err = self.rpcClient.CallContext(ctx, &json, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if json == nil {
		return nil, false, ether.NotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, false, fmt.Errorf("server returned transaction without signature")
	}
	return json, json.BlockNumber == nil, nil
}

func (self *Blockchain) TxStatus(hash ethereum.Hash) (string, uint64, error) {
	option := context.Background()
	tx, pending, err := self.TransactionByHash(option, hash)
	if err == nil {
		// tx exist
		if pending {
			return "", 0, nil
		} else {
			receipt, err := self.client.TransactionReceipt(option, hash)
			if err != nil {
				log.Println("Get receipt err: ", err.Error())
				log.Printf("Receipt: %+v", receipt)
				if receipt != nil {
					// only byzantium has status field at the moment
					// mainnet, ropsten are byzantium, other chains such as
					// devchain, kovan are not
					if self.chainType == "byzantium" {
						if receipt.Status == 1 {
							// successful tx
							return "mined", tx.BlockNumber().Uint64(), nil
						} else {
							// failed tx
							return "failed", tx.BlockNumber().Uint64(), nil
						}
					} else {
						return "mined", tx.BlockNumber().Uint64(), nil
					}
				} else {
					// networking issue
					return "", 0, err
				}
			} else {
				if receipt.Status == 1 {
					// successful tx
					return "mined", tx.BlockNumber().Uint64(), nil
				} else {
					// failed tx
					return "failed", tx.BlockNumber().Uint64(), nil
				}
			}
		}
	} else {
		if err == ether.NotFound {
			// tx doesn't exist. it failed
			return "lost", 0, nil
		} else {
			// networking issue
			return "", 0, err
		}
	}
}

func NewBlockchain(
	client *rpc.Client,
	etherCli *ethclient.Client,
	intermediateSigner Signer, nonceIntermediate NonceCorpus,
	chainType string) (*Blockchain, error) {
	log.Printf("intermediate address: %s", intermediateSigner.GetAddress().Hex())
	return &Blockchain{
		rpcClient:          client,
		client:             etherCli,
		intermediateSigner: intermediateSigner,
		tokens:             []common.Token{},
		nonceIntermediate:  nonceIntermediate,
	}, nil
}
