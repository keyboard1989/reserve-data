package exchange

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	HUOBI_EPSILON float64 = 0.0000000001 // 10e-10
	EXPIRED       uint64  = 3600000      //1 hour in mili second
)

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

type Huobi struct {
	interf            HuobiInterface
	pairs             []common.TokenPair
	addresses         *common.ExchangeAddresses
	exchangeInfo      *common.ExchangeInfo
	fees              common.ExchangeFees
	blockchain        Blockchain
	intermediatorAddr ethereum.Address
	storage           Storage
}

func (self *Huobi) MarshalText() (text []byte, err error) {
	return []byte(self.ID()), nil
}

func (self *Huobi) TokenAddresses() map[string]ethereum.Address {
	return self.addresses.GetData()
}

func (self *Huobi) Address(token common.Token) (ethereum.Address, bool) {

	_, supported := self.addresses.Get(token.ID)
	addr := self.intermediatorAddr
	return addr, supported
}

func (self *Huobi) UpdateAllDepositAddresses(address string, timepoint uint64) {
	data := self.addresses.GetData()
	for k, _ := range data {
		self.addresses.Update(k, ethereum.HexToAddress(address))
	}
}

func (self *Huobi) UpdateDepositAddress(token common.Token, address string) {
	liveAddress, _ := self.interf.GetDepositAddress(strings.ToLower(token.ID))
	if liveAddress.Address != "" {
		self.addresses.Update(token.ID, ethereum.HexToAddress(liveAddress.Address))
	} else {
		self.addresses.Update(token.ID, ethereum.HexToAddress(address))
	}
}

func (self *Huobi) UpdatePrecisionLimit(pair common.TokenPair, symbols HuobiExchangeInfo) {
	pairName := strings.ToLower(pair.Base.ID) + strings.ToLower(pair.Quote.ID)
	for _, symbol := range symbols.Data {
		if symbol.Base+symbol.Quote == pairName {
			exchangePrecisionLimit := common.ExchangePrecisionLimit{}
			exchangePrecisionLimit.Precision.Amount = symbol.AmountPrecision
			exchangePrecisionLimit.Precision.Price = symbol.PricePrecision
			self.exchangeInfo.Update(pair.PairID(), exchangePrecisionLimit)
			break
		}
	}
}

func (self *Huobi) UpdatePairsPrecision() {
	exchangeInfo, err := self.interf.GetExchangeInfo()
	if err != nil {
		log.Printf("Get exchange info failed: %s\n", err)
	} else {
		for _, pair := range self.pairs {
			self.UpdatePrecisionLimit(pair, exchangeInfo)
		}
	}
}

func (self *Huobi) GetInfo() (common.ExchangeInfo, error) {
	return *self.exchangeInfo, nil
}

func (self *Huobi) GetExchangeInfo(pair common.TokenPairID) (common.ExchangePrecisionLimit, error) {
	data, err := self.exchangeInfo.Get(pair)
	return data, err
}

func (self *Huobi) GetFee() common.ExchangeFees {
	return self.fees
}

func (self *Huobi) ID() common.ExchangeID {
	return common.ExchangeID("huobi")
}

func (self *Huobi) TokenPairs() []common.TokenPair {
	return self.pairs
}

func (self *Huobi) Name() string {
	return "huobi"
}

func (self *Huobi) QueryOrder(symbol string, id uint64) (done float64, remaining float64, finished bool, err error) {
	result, err := self.interf.OrderStatus(symbol, id)
	if err != nil {
		return 0, 0, false, err
	} else {
		done, _ := strconv.ParseFloat(result.Data.ExecutedQty, 64)
		total, _ := strconv.ParseFloat(result.Data.OrigQty, 64)
		return done, total - done, total-done < HUOBI_EPSILON, nil
	}
}

func (self *Huobi) Trade(tradeType string, base common.Token, quote common.Token, rate float64, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	result, err := self.interf.Trade(tradeType, base, quote, rate, amount, timepoint)
	symbol := base.ID + quote.ID

	if err != nil {
		return "", 0, 0, false, err
	} else {
		orderID, _ := strconv.ParseUint(result.OrderID, 10, 64)
		done, remaining, finished, err := self.QueryOrder(
			base.ID+quote.ID,
			orderID,
		)
		id := fmt.Sprintf("%s_%s", result.OrderID, symbol)
		return id, done, remaining, finished, err
	}
}

func (self *Huobi) Withdraw(token common.Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	withdrawID, err := self.interf.Withdraw(token, amount, address, timepoint)
	if err != nil {
		return "", err
	}
	// this magical logic base on inspection on huobi website
	result := withdrawID + "01"
	return result, err
}

func (self *Huobi) CancelOrder(id common.ActivityID) error {
	idParts := strings.Split(id.EID, "_")
	idNo, err := strconv.ParseUint(idParts[0], 10, 64)
	if err != nil {
		return err
	}
	symbol := idParts[1]
	result, err := self.interf.CancelOrder(symbol, idNo)
	if err != nil {
		return err
	}
	if result.Status != "ok" {
		return errors.New("Couldn't cancel order id " + id.EID)
	}
	return nil
}

func (self *Huobi) FetchOnePairData(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	defer wg.Done()
	result := common.ExchangePrice{}

	timestamp := common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Timestamp = timestamp
	result.Valid = true
	resp_data, err := self.interf.GetDepthOnePair(pair)
	returnTime := common.GetTimestamp()
	result.ReturnTime = returnTime
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		if resp_data.Status != "ok" {
			result.Valid = false
		} else {
			for _, buy := range resp_data.Tick.Bids {
				quantity := buy[1]
				rate := buy[0]
				result.Bids = append(
					result.Bids,
					common.PriceEntry{
						quantity,
						rate,
					},
				)
			}
			for _, sell := range resp_data.Tick.Asks {
				quantity := sell[1]
				rate := sell[0]
				result.Asks = append(
					result.Asks,
					common.PriceEntry{
						quantity,
						rate,
					},
				)
			}
		}
	}
	data.Store(pair.PairID(), result)
}

func (self *Huobi) FetchPriceData(timepoint uint64) (map[common.TokenPairID]common.ExchangePrice, error) {
	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs := self.pairs
	for _, pair := range pairs {
		wait.Add(1)
		go self.FetchOnePairData(&wait, pair, &data, timepoint)
	}
	wait.Wait()
	result := map[common.TokenPairID]common.ExchangePrice{}
	data.Range(func(key, value interface{}) bool {
		result[key.(common.TokenPairID)] = value.(common.ExchangePrice)
		return true
	})
	return result, nil
}

func (self *Huobi) OpenOrdersForOnePair(
	wg *sync.WaitGroup,
	pair common.TokenPair,
	data *sync.Map,
	timepoint uint64) {

	// defer wg.Done()

	// result, err := self.interf.OpenOrdersForOnePair(pair, timepoint)

	//TODO: complete open orders for one pair
}

func (self *Huobi) FetchOrderData(timepoint uint64) (common.OrderEntry, error) {
	result := common.OrderEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	result.Data = []common.Order{}

	wait := sync.WaitGroup{}
	data := sync.Map{}
	pairs := self.pairs
	for _, pair := range pairs {
		wait.Add(1)
		go self.OpenOrdersForOnePair(&wait, pair, &data, timepoint)
	}
	wait.Wait()

	result.ReturnTime = common.GetTimestamp()

	data.Range(func(key, value interface{}) bool {
		orders := value.([]common.Order)
		result.Data = append(result.Data, orders...)
		return true
	})
	return result, nil
}

func (self *Huobi) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	resp_data, err := self.interf.GetInfo()
	result.ReturnTime = common.GetTimestamp()
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	} else {
		result.AvailableBalance = map[string]float64{}
		result.LockedBalance = map[string]float64{}
		result.DepositBalance = map[string]float64{}
		if resp_data.Status != "ok" {
			result.Valid = false
			result.Error = fmt.Sprintf("Cannot fetch ebalance")
		} else {
			balances := resp_data.Data.List
			for _, b := range balances {
				tokenID := strings.ToUpper(b.Currency)
				_, exist := common.SupportedTokens[tokenID]
				if exist {
					balance, _ := strconv.ParseFloat(b.Balance, 64)
					if b.Type == "trade" {
						result.AvailableBalance[tokenID] = balance
					} else {
						result.LockedBalance[tokenID] = balance
					}
					result.DepositBalance[tokenID] = 0
				}
			}
			return result, nil
		}
	}
	return result, nil
}

func (self *Huobi) FetchOnePairTradeHistory(
	wait *sync.WaitGroup,
	data *sync.Map,
	pair common.TokenPair) {

	defer wait.Done()
	result := []common.TradeHistory{}
	resp, err := self.interf.GetAccountTradeHistory(pair.Base, pair.Quote)
	if err != nil {
		log.Printf("Cannot fetch data for pair %s%s: %s", pair.Base.ID, pair.Quote.ID, err.Error())
	}
	pairString := pair.PairID()
	for _, trade := range resp.Data {
		price, _ := strconv.ParseFloat(trade.Price, 64)
		quantity, _ := strconv.ParseFloat(trade.Amount, 64)
		historyType := "sell"
		if trade.Type == "buy-limit" {
			historyType = "buy"
		}
		tradeHistory := common.TradeHistory{
			strconv.FormatUint(trade.ID, 10),
			price,
			quantity,
			historyType,
			trade.Timestamp,
		}
		result = append(result, tradeHistory)
	}
	data.Store(pairString, result)
}

func (self *Huobi) FetchTradeHistory(timepoint uint64) (map[common.TokenPairID][]common.TradeHistory, error) {
	result := map[common.TokenPairID][]common.TradeHistory{}
	data := sync.Map{}
	pairs := self.pairs
	wait := sync.WaitGroup{}
	for _, pair := range pairs {
		wait.Add(1)
		go self.FetchOnePairTradeHistory(&wait, &data, pair)
	}
	wait.Wait()
	data.Range(func(key, value interface{}) bool {
		result[key.(common.TokenPairID)] = value.([]common.TradeHistory)
		return true
	})
	return result, nil
}

func getDepositInfo(id common.ActivityID) (string, float64, string) {
	idParts := strings.Split(id.EID, "|")
	txID := idParts[0]
	sentAmount, err := strconv.ParseFloat(idParts[2], 64)
	if err != nil {
		log.Println("The ID is malform, cannot get Amount from EID")
	}
	tokenID := idParts[1]
	return txID, sentAmount, tokenID
}

func (self *Huobi) Send2ndTransaction(amount float64, token common.Token, exchangeAddress ethereum.Address) (*types.Transaction, error) {
	IAmount := getBigIntFromFloat(amount, token.Decimal)
	// Check balance, removed from huobi's blockchain object.
	// currBalance := self.blockchain.CheckBalance(token)
	// log.Printf("current balance of token %s is %d", token.ID, currBalance)
	// //self.blockchain.
	// if currBalance.Cmp(IAmount) < 0 {
	// 	log.Printf("balance is not enough, wait till next check")
	// 	return nil, errors.New("balance is not enough")
	// }
	var tx *types.Transaction
	var err error
	if token.ID == "ETH" {
		tx, err = self.blockchain.SendETHFromAccountToExchange(IAmount, exchangeAddress)
	} else {
		tx, err = self.blockchain.SendTokenFromAccountToExchange(IAmount, exchangeAddress, ethereum.HexToAddress(token.Address))
	}
	if err != nil {
		log.Printf("ERROR: Can not send transaction to exchange: %v", err)
		return nil, err
	}
	log.Printf("Transaction submitted. Tx is: \n %v", tx)
	return tx, nil

}

func (self *Huobi) PendingIntermediateTxs(timepoint uint64) (map[common.ActivityID]common.TXEntry, error) {
	result := make(map[common.ActivityID]common.TXEntry)
	pendings, err := self.storage.GetPendingIntermediateTXID(timepoint)
	if err != nil {
		return nil, err
	}
	for _, actID := range pendings {
		tx2, err := self.storage.GetIntermedatorTx(actID)
		if err == nil {
			result[actID] = tx2
		} else {
			log.Printf("the current acitivy doesn't have the 2nd tx yet: ", err)
		}
	}
	return result, nil
}

func (self *Huobi) DepositStatus(id common.ActivityID, timepoint uint64) (string, error) {
	txID, sentAmount, tokenID := getDepositInfo(id)
	tx2Entry, ok := self.storage.GetIntermedatorTx(id)
	err := self.storage.StorePendingIntermediateTxID(common.GetTimepoint(), id)
	if err != nil {
		return "", err
	}
	if ok != nil {
		//if the 2nd transaction is not in the current Deposit status, check the 1st tx first.
		status, blockno, err := self.blockchain.TxStatus(ethereum.HexToHash(txID))
		if err != nil {
			log.Println("Can not get TX status")
		}
		log.Printf("Status was %s at block %d ", status, blockno)
		if status == "mined" {
			//if it is mined, send 2nd tx.
			log.Printf("found a new deposit status, which deposit %.5f %s. Procceed to send it to Huobi", sentAmount, tokenID)
			//check if the token is supported
			token, err := common.GetToken(tokenID)
			if err != nil {
				log.Printf(" Token %s is not supported", tokenID)
				return "", err
			}
			exchangeAddress, ok := self.addresses.Get(tokenID)
			if !ok {
				log.Print(" Wrong token address configuration ")
				return "", errors.New("Wrong token address configuration")
			}
			tx2, err := self.Send2ndTransaction(sentAmount, token, exchangeAddress)
			if err != nil {
				log.Println("can't send 2nd transaction to the exchange")
				return "failed", err
			}
			Txhash := tx2.Hash().Hex()
			err = self.storage.StoreIntermediateTx(Txhash, self.Name(), tokenID, "submitted", "", sentAmount, common.GetTimestamp(), id)
			if err != nil {
				log.Printf("Can not store the activity's 2nd transaction")
			}
		} else {
			return "", nil
		}
	} else {
		//if the 2nd transaction is in the Deposit Status, check its status.
		status, _, err := self.blockchain.TxStatus(ethereum.HexToHash(tx2Entry.Hash))
		if err != nil {
			log.Printf("Can not get TX status: %v", err)
			return "", nil
		}
		if status == "mined" {
			log.Println("2nd Transaction is mined. Processed to store it and check the Deposit history")
			err = self.storage.StoreIntermediateTx(tx2Entry.Hash, self.Name(), tokenID, "mined", "", sentAmount, common.GetTimestamp(), id)
			if err != nil {
				log.Printf("Can not store the activity's 2nd transaction ")
			}
			deposits, err := self.interf.DepositHistory()
			if err != nil && deposits.Status != "ok" {
				return "", err
			}
			for _, deposit := range deposits.Data {
				log.Printf("deposit tx is %s, with token %s", deposit.TxHash, deposit.Currency)
				if deposit.TxHash == tx2Entry.Hash {
					if deposit.State == "safe" {
						err = self.storage.StoreIntermediateTx(tx2Entry.Hash, self.Name(), tokenID, "mined", "done", sentAmount, common.GetTimestamp(), id)
						if err != nil {
							log.Printf("Can not store the activity's 2nd transaction ")
						}
						return "done", nil
					}
				}
			}
			return "", errors.New("Deposit doesn't exist. This shouldn't happen unless tx returned from huobi and tx2 are not consistently designed")
		} else if status == "failed" || status == "lost" {
			err = self.storage.StoreIntermediateTx(tx2Entry.Hash, self.Name(), tokenID, "failed", "failed", sentAmount, common.GetTimestamp(), id)
			return "failed", nil
		}
	}
	return "", nil
}

func (self *Huobi) WithdrawStatus(id common.ActivityID, timepoint uint64) (string, string, error) {
	withdrawID, _ := strconv.ParseUint(id.EID, 10, 64)
	withdraws, err := self.interf.WithdrawHistory()
	if err != nil {
		return "", "", nil
	}
	log.Printf("Withdrawal id: %s", withdrawID)
	for _, withdraw := range withdraws.Data {
		if withdraw.TxID == withdrawID {
			if withdraw.State == "confirmed" {
				return "done", withdraw.TxHash, nil
			}
			return "", withdraw.TxHash, nil
		}
	}
	return "", "", errors.New("Withdrawal doesn't exist. This shouldn't happen unless tx returned from withdrawal from huobi and activity ID are not consistently designed")
}

func (self *Huobi) OrderStatus(id common.ActivityID, timepoint uint64) (string, error) {
	tradeID := id.EID
	parts := strings.Split(tradeID, "_")
	orderID, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		// if this crashes, it means core put malformed activity ID
		panic(err)
	}
	symbol := parts[1]
	order, err := self.interf.OrderStatus(symbol, orderID)
	if err != nil {
		return "", err
	}
	if order.Data.State == "pre-submitted" || order.Data.State == "submitting" || order.Data.State == "submitted" || order.Data.State == "partial-filled" || order.Data.State == "partial-canceled" {
		return "", nil
	} else {
		return "done", nil
	}
}

func NewHuobi(addressConfig map[string]string, feeConfig common.ExchangeFees, interf HuobiInterface,
	intorSigner Signer, ethEndpoint string, intorAddr ethereum.Address, storage Storage) *Huobi {
	pairs, fees := getExchangePairsAndFeesFromConfig(addressConfig, feeConfig, "huobi")
	bc, err := NewBlockchain(intorSigner, ethEndpoint)
	if err != nil {
		log.Printf("Cant create Huobi's blockchain: %v", err)
		panic(err)
	}
	return &Huobi{
		interf,
		pairs,
		common.NewExchangeAddresses(),
		common.NewExchangeInfo(),
		fees,
		*bc,
		intorAddr,
		storage,
	}
}
