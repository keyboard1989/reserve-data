package fetcher

import (
	"errors"
	"log"
	"sync"

	"github.com/KyberNetwork/reserve-data/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Fetcher struct {
	storage    Storage
	exchanges  []Exchange
	blockchain Blockchain
	runner     FetcherRunner
	rmaddr     ethereum.Address
}

func NewFetcher(
	storage Storage,
	runner FetcherRunner,
	address ethereum.Address) *Fetcher {
	return &Fetcher{
		storage:    storage,
		exchanges:  []Exchange{},
		blockchain: nil,
		runner:     runner,
		rmaddr:     address,
	}
}

func (self *Fetcher) getExchange(id string) Exchange {
	for _, ex := range self.exchanges {
		if string(ex.ID()) == id {
			return ex
		}
	}
	panic(errors.New("Exchange " + id + " is not registered in fetcher"))
}

func (self *Fetcher) SetBlockchain(blockchain Blockchain) {
	self.blockchain = blockchain
}

func (self *Fetcher) AddExchange(exchange Exchange) {
	self.exchanges = append(self.exchanges, exchange)
}

func (self *Fetcher) Stop() error {
	return self.runner.Stop()
}

func (self *Fetcher) Run() error {
	log.Printf("Fetcher runner is starting...")
	self.runner.Start()
	log.Printf("Fetcher runner is running...")
	go self.RunOrderbookFetcher()
	go self.RunAuthDataFetcher()
	return nil
}

func (self *Fetcher) RunAuthDataFetcher() {
	for {
		log.Printf("waiting for signal from runner auth data channel")
		t := <-self.runner.GetAuthDataTicker()
		log.Printf("got signal in auth data channel with timestamp %d", common.TimeToTimepoint(t))
		self.FetchAllAuthData(common.TimeToTimepoint(t))
		log.Printf("fetched data from exchanges")
	}
}

func (self *Fetcher) FetchAllAuthData(timepoint uint64) {
	snapshot := common.AuthDataSnapshot{
		Valid:              true,
		Timestamp:          common.GetTimestamp(),
		ExchangeBalances:   map[common.ExchangeID]common.EBalanceEntry{},
		BlockchainBalances: map[string]common.BalanceEntry{},
		PendingActivities:  []common.ActivityRecord{},
	}
	bbalances := map[string]common.BalanceEntry{}
	ebalances := sync.Map{}
	estatuses := sync.Map{}
	bstatuses := sync.Map{}
	pendings, err := self.storage.GetPendingActivities()
	if err != nil {
		log.Printf("Getting pending activites failed: %s\n", err)
		return
	}
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.FetchAuthDataFromExchange(
			&wait, exchange, &ebalances, &estatuses,
			pendings, timepoint)
	}
	wait.Wait()
	self.FetchAuthDataFromBlockchain(
		bbalances, &bstatuses, pendings, timepoint)
	snapshot.ReturnTime = common.GetTimestamp()
	err = self.PersistSnapshot(
		&ebalances, bbalances, &estatuses, &bstatuses,
		pendings, &snapshot, timepoint)
	if err != nil {
		log.Printf("Storing exchange balances failed: %s\n", err)
		return
	}
}

func (self *Fetcher) FetchAuthDataFromBlockchain(
	allBalances map[string]common.BalanceEntry,
	allStatuses *sync.Map,
	pendings []common.ActivityRecord,
	timepoint uint64) {
	// we apply double check strategy to mitigate race condition on exchange side like this:
	// 1. Get list of pending activity status (A)
	// 2. Get list of balances (B)
	// 3. Get list of pending activity status again (C)
	// 4. if C != A, repeat 1, otherwise return A, B
	var balances map[string]common.BalanceEntry
	var statuses map[common.ActivityID]common.ActivityStatus
	var err error
	for {
		preStatuses := self.FetchStatusFromBlockchain(pendings)
		balances, err = self.FetchBalanceFromBlockchain(timepoint)
		if err != nil {
			log.Printf("Fetching blockchain balances failed: %v\n", err)
			break
		}
		statuses = self.FetchStatusFromBlockchain(pendings)
		if unchanged(preStatuses, statuses) {
			break
		}
	}
	if err == nil {
		for k, v := range balances {
			allBalances[k] = v
		}
		for id, activityStatus := range statuses {
			allStatuses.Store(id, activityStatus)
		}
	}
}

func (self *Fetcher) FetchBalanceFromBlockchain(timepoint uint64) (map[string]common.BalanceEntry, error) {
	return self.blockchain.FetchBalanceData(self.rmaddr, timepoint)
}

func (self *Fetcher) FetchStatusFromBlockchain(pendings []common.ActivityRecord) map[common.ActivityID]common.ActivityStatus {
	result := map[common.ActivityID]common.ActivityStatus{}
	for _, activity := range pendings {
		if activity.MiningStatus != "mined" && activity.MiningStatus != "failed" {
			switch activity.Action {
			case "set_rates", "deposit", "withdraw":
				tx := ethereum.HexToHash(activity.Result["tx"].(string))
				if !tx.Big().IsInt64() || tx.Big().Int64() != 0 {
					isMined, err := self.blockchain.IsMined(tx)
					if isMined {
						result[activity.ID] = common.ActivityStatus{
							activity.ExchangeStatus,
							activity.Result["tx"].(string),
							"mined",
							err,
						}
					}
				}
			}
		}
	}
	return result
}

func unchanged(pre, post map[common.ActivityID]common.ActivityStatus) bool {
	if len(pre) != len(post) {
		return false
	} else {
		for k, v := range pre {
			vpost, found := post[k]
			if !found {
				return false
			}
			if v.ExchangeStatus != vpost.ExchangeStatus || v.MiningStatus != vpost.MiningStatus {
				return false
			}
		}
	}
	return true
}

func (self *Fetcher) PersistSnapshot(
	ebalances *sync.Map,
	bbalances map[string]common.BalanceEntry,
	estatuses *sync.Map,
	bstatuses *sync.Map,
	pendings []common.ActivityRecord,
	snapshot *common.AuthDataSnapshot,
	timepoint uint64) error {

	allEBalances := map[common.ExchangeID]common.EBalanceEntry{}
	ebalances.Range(func(key, value interface{}) bool {
		v := value.(common.EBalanceEntry)
		allEBalances[key.(common.ExchangeID)] = v
		if !v.Valid {
			snapshot.Valid = false
			snapshot.Error = v.Error
		}
		return true
	})

	pendingActivities := []common.ActivityRecord{}
	for _, activity := range pendings {
		status, _ := estatuses.Load(activity.ID)
		activityStatus := status.(common.ActivityStatus)
		if activityStatus.Error == nil {
			activity.ExchangeStatus = activityStatus.ExchangeStatus
			activity.Result["tx"] = activityStatus.Tx
		} else {
			snapshot.Valid = false
			snapshot.Error = activityStatus.Error.Error()
		}
		status, _ = bstatuses.Load(activity.ID)
		activityStatus = status.(common.ActivityStatus)
		if activityStatus.Error == nil {
			activity.MiningStatus = activityStatus.MiningStatus
		} else {
			snapshot.Valid = false
			snapshot.Error = activityStatus.Error.Error()
		}
		if activity.IsPending() {
			pendingActivities = append(pendingActivities, activity)
		} else {
			err := self.storage.UpdateActivity(activity.ID, activity)
			if err != nil {
				snapshot.Valid = false
				snapshot.Error = err.Error()
			}
		}
	}
	snapshot.ExchangeBalances = allEBalances
	snapshot.BlockchainBalances = bbalances
	snapshot.PendingActivities = pendingActivities
	return self.storage.StoreAuthSnapshot(snapshot, timepoint)
}

func (self *Fetcher) FetchAuthDataFromExchange(
	wg *sync.WaitGroup, exchange Exchange,
	allBalances *sync.Map, allStatuses *sync.Map,
	pendings []common.ActivityRecord,
	timepoint uint64) {
	defer wg.Done()
	// we apply double check strategy to mitigate race condition on exchange side like this:
	// 1. Get list of pending activity status (A)
	// 2. Get list of balances (B)
	// 3. Get list of pending activity status again (C)
	// 4. if C != A, repeat 1, otherwise return A, B
	var balances common.EBalanceEntry
	var statuses map[common.ActivityID]common.ActivityStatus
	var err error
	for {
		preStatuses := self.FetchStatusFromExchange(exchange, pendings, timepoint)
		balances, err = exchange.FetchEBalanceData(timepoint)
		if err != nil {
			log.Printf("Fetching exchange balances from %s failed: %v\n", exchange.Name(), err)
			break
		}
		statuses = self.FetchStatusFromExchange(exchange, pendings, timepoint)
		if unchanged(preStatuses, statuses) {
			break
		}
	}
	if err == nil {
		allBalances.Store(exchange.ID(), balances)
		for id, activityStatus := range statuses {
			allStatuses.Store(id, activityStatus)
		}
	}
}

func (self *Fetcher) FetchStatusFromExchange(exchange Exchange, pendings []common.ActivityRecord, timepoint uint64) map[common.ActivityID]common.ActivityStatus {
	result := map[common.ActivityID]common.ActivityStatus{}
	for _, activity := range pendings {
		if activity.Destination == string(exchange.ID()) && activity.ExchangeStatus != "done" && activity.ExchangeStatus != "failed" {
			var err error
			var status string
			var tx string
			id := activity.ID
			if activity.Action == "trade" {
				status, err = exchange.OrderStatus(id, timepoint)
			} else if activity.Action == "deposit" {
				status, err = exchange.DepositStatus(id, timepoint)
			} else if activity.Action == "withdraw" {
				status, tx, err = exchange.WithdrawStatus(id, timepoint)
			} else {
				continue
			}
			result[id] = common.ActivityStatus{
				status,
				tx,
				activity.MiningStatus,
				err,
			}
		}
	}
	return result
}

func (self *Fetcher) RunOrderbookFetcher() {
	for {
		log.Printf("waiting for signal from runner orderbook channel")
		t := <-self.runner.GetOrderbookTicker()
		log.Printf("got signal in orderbook channel with timestamp %d", common.TimeToTimepoint(t))
		self.FetchOrderbook(common.TimeToTimepoint(t))
		log.Printf("fetched data from exchanges")
	}
}

func (self *Fetcher) FetchOrderbook(timepoint uint64) {
	data := NewConcurrentAllPriceData()
	// start fetching
	wait := sync.WaitGroup{}
	for _, exchange := range self.exchanges {
		wait.Add(1)
		go self.fetchPriceFromExchange(&wait, exchange, data, timepoint)
	}
	wait.Wait()
	err := self.storage.StorePrice(data.GetData(), timepoint)
	if err != nil {
		log.Printf("Storing data failed: %s\n", err)
	}
}

func (self *Fetcher) fetchPriceFromExchange(wg *sync.WaitGroup, exchange Exchange, data *ConcurrentAllPriceData, timepoint uint64) {
	defer wg.Done()
	exdata, err := exchange.FetchPriceData(timepoint)
	if err != nil {
		log.Printf("Fetching data from %s failed: %v\n", exchange.Name(), err)
	}
	for pair, exchangeData := range exdata {
		data.SetOnePrice(exchange.ID(), pair, exchangeData)
	}
}

// func (self *Fetcher) fetchingFromExchanges() {
// 	for {
// 		log.Printf("waiting for signal from runner for exchange ticker")
// 		t := <-self.runner.GetExchangeTicker()
// 		log.Printf("got signal in exchange ticker with timestamp %d", common.TimeToTimepoint(t))
// 		self.fetchAllFromExchanges(common.TimeToTimepoint(t))
// 		log.Printf("fetched data from exchanges")
// 	}
// }
//
// func (self *Fetcher) fetchingFromBlockchain() {
// 	for {
// 		t := <-self.runner.GetBlockchainTicker()
// 		self.fetchAllFromBlockchain(common.TimeToTimepoint(t))
// 	}
// }
//
//
// func (self *Fetcher) fetchEBalanceFromExchange(wg *sync.WaitGroup, exchange Exchange, data *sync.Map, timepoint uint64) {
// 	defer wg.Done()
// 	exdata, err := exchange.FetchEBalanceData(timepoint)
// 	if err != nil {
// 		log.Printf("Fetching exchange balances from %s failed: %v\n", exchange.Name(), err)
// 	}
// 	log.Printf("Fetched balance data from %s: %v", exchange.ID(), exdata)
// 	data.Store(exchange.ID(), exdata)
// }
//
//
// func (self *Fetcher) fetchAllPrices(w *sync.WaitGroup, timepoint uint64) {
// 	defer w.Done()
// 	data := NewConcurrentAllPriceData()
// 	// start fetching
// 	wait := sync.WaitGroup{}
// 	for _, exchange := range self.exchanges {
// 		wait.Add(1)
// 		go self.fetchPriceFromExchange(&wait, exchange, data, timepoint)
// 	}
// 	wait.Wait()
// 	err := self.storage.StorePrice(data.GetData(), timepoint)
// 	if err != nil {
// 		log.Printf("Storing data failed: %s\n", err)
// 	}
// }
//
// func (self *Fetcher) fetchAllEBalances(w *sync.WaitGroup, timepoint uint64) {
// 	defer w.Done()
// 	data := sync.Map{}
// 	// start fetching
// 	wait := sync.WaitGroup{}
// 	for _, exchange := range self.exchanges {
// 		wait.Add(1)
// 		go self.fetchEBalanceFromExchange(&wait, exchange, &data, timepoint)
// 	}
// 	wait.Wait()
// 	ebalances := map[common.ExchangeID]common.EBalanceEntry{}
// 	data.Range(func(key, value interface{}) bool {
// 		ebalances[key.(common.ExchangeID)] = value.(common.EBalanceEntry)
// 		return true
// 	})
// 	err := self.storage.StoreEBalance(ebalances, timepoint)
// 	if err != nil {
// 		log.Printf("Storing exchange balances failed: %s\n", err)
// 	}
// }
//
// func (self *Fetcher) fetchAllBalances(w *sync.WaitGroup, timepoint uint64) {
// 	defer w.Done()
// 	data, err := self.blockchain.FetchBalanceData(self.rmaddr, timepoint)
// 	if err != nil {
// 		log.Printf("Fetching data from blockchain failed: %s\n", err)
// 	}
// 	err = self.storage.StoreBalance(data, timepoint)
// 	// fmt.Printf("balance data: %v\n", data)
// 	if err != nil {
// 		log.Printf("Storing balance data failed: %s\n", err)
// 	}
// }
//
// func (self *Fetcher) fetchStatusFromBlockchainAndUpdate(w *sync.WaitGroup, activity common.ActivityRecord, timepoint uint64) {
// 	defer w.Done()
// 	action := activity.Action
// 	id := activity.ID
// 	destination := activity.Destination
// 	tx := ethereum.HexToHash(activity.Result["tx"].(string))
// 	isMined, err := self.blockchain.IsMined(tx)
// 	if err != nil {
// 		log.Printf("Fetching Tx mining status failed: %s", err.Error())
// 	}
// 	if isMined {
// 		self.storage.UpdateActivityStatus(action, id, destination, "mined")
// 	}
// }
//
// func (self *Fetcher) fetchStatusFromExchangeAndUpdate(w *sync.WaitGroup, activity common.ActivityRecord, timepoint uint64) {
// 	defer w.Done()
// 	action := activity.Action
// 	id := activity.ID
// 	destination := activity.Destination
// 	exchange := self.getExchange(destination)
// 	var err error
// 	var status string
// 	if activity.Action == "trade" {
// 		status, err = exchange.OrderStatus(id, timepoint)
// 	} else if activity.Action == "deposit" {
// 		status, err = exchange.DepositStatus(id, timepoint)
// 	} else if activity.Action == "withdraw" {
// 		status, err = exchange.WithdrawStatus(id, timepoint)
// 	}
// 	if err == nil && status != activity.Status && status != "" {
// 		self.storage.UpdateActivityStatus(action, id, destination, status)
// 	}
// }
//
// func (self *Fetcher) fetchActivityStatus(w *sync.WaitGroup, timepoint uint64) {
// 	defer w.Done()
// 	wg := sync.WaitGroup{}
// 	pendings, err := self.storage.GetPendingActivities()
// 	if err != nil {
// 		log.Printf("Getting pending activities from storage failed: %s", err.Error())
// 	}
// 	for _, activity := range pendings {
// 		if activity.Action == "set_rates" {
// 			wg.Add(1)
// 			go self.fetchStatusFromBlockchainAndUpdate(&wg, activity, timepoint)
// 		} else {
// 			wg.Add(1)
// 			go self.fetchStatusFromExchangeAndUpdate(&wg, activity, timepoint)
// 		}
// 	}
// 	wg.Wait()
// }
//
// func (self *Fetcher) fetchAllRates(w *sync.WaitGroup, timepoint uint64) {
// 	defer w.Done()
// 	log.Printf("Fetching all rates from blockchain...")
// 	sources := []common.Token{}
// 	dests := []common.Token{}
// 	pairs := map[common.TokenPairID]bool{}
// 	for _, ex := range self.exchanges {
// 		tokenPairs := ex.TokenPairs()
// 		for _, p := range tokenPairs {
// 			_, exist := pairs[p.PairID()]
// 			if !exist {
// 				pairs[p.PairID()] = true
// 				sources = append(sources, p.Base)
// 				dests = append(dests, p.Quote)
// 			}
// 		}
// 	}
// 	data, err := self.blockchain.FetchRates(sources, dests, timepoint)
// 	if err != nil {
// 		log.Printf("Fetching rate data from blockchain failed: %s\n", err)
// 	}
// 	err = self.storage.StoreRate(data, timepoint)
// 	// fmt.Printf("balance data: %v\n", data)
// 	if err != nil {
// 		log.Printf("Storing balance data failed: %s\n", err)
// 	}
// }
//
// func (self *Fetcher) fetchAllFromExchanges(timepoint uint64) {
// 	log.Printf("Fetching all data from exchanges...")
// 	wait := sync.WaitGroup{}
// 	wait.Add(1)
// 	go self.fetchAllPrices(&wait, timepoint)
// 	wait.Add(1)
// 	go self.fetchAllEBalances(&wait, timepoint)
// 	wait.Add(1)
// 	go self.fetchActivityStatus(&wait, timepoint)
// 	log.Printf("Waiting price, balance, order data from exchanges...")
// 	wait.Wait()
// }
//
// func (self *Fetcher) fetchAllFromBlockchain(timepoint uint64) {
// 	log.Printf("Fetching data from blockchain...")
// 	wait := sync.WaitGroup{}
// 	wait.Add(1)
// 	go self.fetchAllBalances(&wait, timepoint)
// 	wait.Add(1)
// 	go self.fetchAllRates(&wait, timepoint)
// 	wait.Wait()
// }
