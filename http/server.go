package http

import (
	"fmt"
	"github.com/KyberNetwork/reserve-data"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app  reserve.ReserveData
	core reserve.ReserveCore
	host string
	r    *gin.Engine
}

const MAX_TIMESPOT uint64 = 18446744073709551615

func getTimePoint(c *gin.Context) uint64 {
	timestamp := c.DefaultQuery("timestamp", "")
	if timestamp == "" {
		fmt.Printf("Interpreted timestamp(%s) to default - %s\n", timestamp, MAX_TIMESPOT)
		return MAX_TIMESPOT
	} else {
		timepoint, err := strconv.ParseUint(timestamp, 10, 64)
		if err != nil {
			fmt.Printf("Interpreted timestamp(%s) to default - %s\n", timestamp, MAX_TIMESPOT)
			return MAX_TIMESPOT
		} else {
			fmt.Printf("Interpreted timestamp(%s) to %s\n", timestamp, timepoint)
			return timepoint
		}
	}
}

func (self *HTTPServer) AllPrices(c *gin.Context) {
	fmt.Printf("Getting all prices \n")
	data, err := self.app.GetAllPrices(getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success":   true,
				"version":   data.Version,
				"timestamp": data.Timestamp,
				"data":      data.Data,
			},
		)
	}
}

func (self *HTTPServer) Price(c *gin.Context) {
	base := c.Param("base")
	quote := c.Param("quote")
	fmt.Printf("Getting price for %s - %s \n", base, quote)
	pair, err := common.NewTokenPair(base, quote)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": "Token pair is not supported"},
		)
	} else {
		data, err := self.app.GetOnePrice(pair.PairID(), getTimePoint(c))
		if err != nil {
			c.JSON(
				http.StatusOK,
				gin.H{"success": false, "reason": err.Error()},
			)
		} else {
			c.JSON(
				http.StatusOK,
				gin.H{
					"success":   true,
					"version":   data.Version,
					"timestamp": data.Timestamp,
					"exchanges": data.Data,
				},
			)
		}
	}
}

func (self *HTTPServer) AllBalances(c *gin.Context) {
	fmt.Printf("Getting all balances \n")
	data, err := self.app.GetAllBalances(getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success":   true,
				"version":   data.Version,
				"timestamp": data.Timestamp,
				"data":      data.Data,
			},
		)
	}
}

func (self *HTTPServer) AllEBalances(c *gin.Context) {
	fmt.Printf("Getting all balances \n")
	data, err := self.app.GetAllEBalances(getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success":   true,
				"version":   data.Version,
				"timestamp": data.Timestamp,
				"data":      data.Data,
			},
		)
	}
}

func (self *HTTPServer) GetRate(c *gin.Context) {
	fmt.Printf("Getting all rates \n")
	data, err := self.app.GetAllRates(getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success":   true,
				"version":   data.Version,
				"timestamp": data.Timestamp,
				"data":      data.Data,
			},
		)
	}
}

func (self *HTTPServer) SetRate(c *gin.Context) {
	sources := c.PostForm("sources")
	dests := c.PostForm("dests")
	rates := c.PostForm("rates")
	blocks := c.PostForm("expiries")
	sourceTokens := []common.Token{}
	for _, source := range strings.Split(sources, "-") {
		token, err := common.GetToken(source)
		if err != nil {
			c.JSON(
				http.StatusOK,
				gin.H{"success": false, "reason": err.Error()},
			)
			return
		} else {
			sourceTokens = append(sourceTokens, token)
		}
	}
	destTokens := []common.Token{}
	for _, dest := range strings.Split(dests, "-") {
		token, err := common.GetToken(dest)
		if err != nil {
			c.JSON(
				http.StatusOK,
				gin.H{"success": false, "reason": err.Error()},
			)
			return
		} else {
			destTokens = append(destTokens, token)
		}
	}
	bigRates := []*big.Int{}
	for _, rate := range strings.Split(rates, "-") {
		r, err := hexutil.DecodeBig(rate)
		if err != nil {
			c.JSON(
				http.StatusOK,
				gin.H{"success": false, "reason": err.Error()},
			)
		} else {
			bigRates = append(bigRates, r)
		}
	}
	expiryBlocks := []*big.Int{}
	for _, expiry := range strings.Split(blocks, "-") {
		r, err := hexutil.DecodeBig(expiry)
		if err != nil {
			c.JSON(
				http.StatusOK,
				gin.H{"success": false, "reason": err.Error()},
			)
		} else {
			expiryBlocks = append(expiryBlocks, r)
		}
	}
	hash, err := self.core.SetRates(sourceTokens, destTokens, bigRates, expiryBlocks)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": true,
				"hash":    hash.Hex(),
			},
		)
	}
}

func (self *HTTPServer) Trade(c *gin.Context) {
	exchangeParam := c.Param("exchangeid")
	baseTokenParam := c.PostForm("base")
	quoteTokenParam := c.PostForm("quote")
	amountParam := c.PostForm("amount")
	rateParam := c.PostForm("rate")
	typeParam := c.PostForm("type")

	fmt.Printf("params: %v\n", c)

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	base, err := common.GetToken(baseTokenParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	quote, err := common.GetToken(quoteTokenParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	rate, err := strconv.ParseFloat(rateParam, 64)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	if typeParam != "sell" && typeParam != "buy" {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": fmt.Sprintf("Trade type of %s is not supported.", typeParam)},
		)
		return
	}
	done, remaining, finished, err := self.core.Trade(
		exchange, typeParam, base, quote, rate, amount, getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success":   true,
			"done":      done,
			"remaining": remaining,
			"finished":  finished,
		},
	)
}

func (self *HTTPServer) Withdraw(c *gin.Context) {
	exchangeParam := c.Param("exchangeid")
	tokenParam := c.PostForm("token")
	amountParam := c.PostForm("amount")

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	token, err := common.GetToken(tokenParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	amount, err := hexutil.DecodeBig(amountParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	fmt.Printf("Withdraw %s %s from %s\n", amount.Text(10), token.ID, exchange.ID())
	err = self.core.Withdraw(exchange, token, amount, getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
		},
	)
}

func (self *HTTPServer) Deposit(c *gin.Context) {
	exchangeParam := c.Param("exchangeid")
	amountParam := c.PostForm("amount")
	tokenParam := c.PostForm("token")

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	token, err := common.GetToken(tokenParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	amount, err := hexutil.DecodeBig(amountParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	fmt.Printf("Depositing %s %s to %s\n", amount.Text(10), token.ID, exchange.ID())
	hash, err := self.core.Deposit(exchange, token, amount, getTimePoint(c))
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"hash":    hash.Hex(),
		},
	)
}

func (self *HTTPServer) GetActivities(c *gin.Context) {
	fmt.Printf("Getting all activity records \n")
	data, err := self.core.GetRecords()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": true,
				"data":    data,
			},
		)
	}
}

func (self *HTTPServer) Run() {
	self.r.GET("/prices", self.AllPrices)
	self.r.GET("/prices/:base/:quote", self.Price)
	self.r.GET("/balances", self.AllBalances)
	self.r.GET("/ebalances", self.AllEBalances)
	self.r.POST("/deposit/:exchangeid", self.Deposit)
	self.r.POST("/withdraw/:exchangeid", self.Withdraw)
	self.r.POST("/trade/:exchangeid", self.Trade)
	self.r.POST("/setrates", self.SetRate)
	self.r.GET("/getrates", self.GetRate)
	self.r.GET("/activities", self.GetActivities)

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Couldn't open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	self.r.Run(self.host)
}

func NewHTTPServer(app reserve.ReserveData, core reserve.ReserveCore, host string) *HTTPServer {
	raven.SetDSN("https://bf15053001464a5195a81bc41b644751:eff41ac715114b20b940010208271b13@sentry.io/228067")

	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))

	return &HTTPServer{
		app, core, host, r,
	}
}
