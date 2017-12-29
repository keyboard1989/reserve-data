package http

import (
	"fmt"
	"github.com/KyberNetwork/reserve-data"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app         reserve.ReserveData
	core        reserve.ReserveCore
	host        string
	authEnabled bool
	auth        Authentication
	r           *gin.Engine
}

const MAX_TIMESPOT uint64 = 18446744073709551615

func getTimePoint(c *gin.Context, useDefault bool) uint64 {
	timestamp := c.DefaultQuery("timestamp", "")
	if timestamp == "" {
		if useDefault {
			log.Printf("Interpreted timestamp to default - %d\n", MAX_TIMESPOT)
			return MAX_TIMESPOT
		} else {
			timepoint := common.GetTimepoint()
			log.Printf("Interpreted timestamp to current time - %d\n", timepoint)
			return uint64(timepoint)
		}
	} else {
		timepoint, err := strconv.ParseUint(timestamp, 10, 64)
		if err != nil {
			log.Printf("Interpreted timestamp(%s) to default - %s\n", timestamp, MAX_TIMESPOT)
			return MAX_TIMESPOT
		} else {
			log.Printf("Interpreted timestamp(%s) to %s\n", timestamp, timepoint)
			return timepoint
		}
	}
}

func IsIntime(nonce string) bool {
	serverTime := common.GetTimepoint()
	nonceInt, err := strconv.ParseInt(nonce, 10, 64)
	if err != nil {
		log.Printf("IsIntime returns false, err: %v", err)
		return false
	}
	difference := nonceInt - int64(serverTime)
	if difference < -10000 || difference > 10000 {
		log.Printf("IsIntime returns false, nonce: %d, serverTime: %d, difference: %d", nonceInt, int64(serverTime), difference)
		return false
	}
	return true
}

// signed message (message = url encoded both query params and post params, keys are sorted) in "signed" header
// using HMAC512
// params must contain "nonce" which is the unixtime in millisecond. The nonce will be invalid
// if it differs from server time more than 10s
func (self *HTTPServer) Authenticated(c *gin.Context, requiredParams []string) (url.Values, bool) {
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "Malformed request package",
			},
		)
		return c.Request.Form, false
	}

	if !self.authEnabled {
		return c.Request.Form, true
	}

	params := c.Request.Form
	if !IsIntime(params.Get("nonce")) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "Your nonce is invalid",
			},
		)
		return c.Request.Form, false
	}

	for _, p := range requiredParams {
		if params.Get(p) == "" {
			c.JSON(
				http.StatusOK,
				gin.H{
					"success": false,
					"reason":  fmt.Sprintf("Required param (%s) is missing. Param name is case sensitive", p),
				},
			)
			return c.Request.Form, false
		}
	}

	signed := c.GetHeader("signed")
	message := c.Request.Form.Encode()
	knsign := self.auth.KNSign(message)
	log.Printf(
		"Signing message(%s) to check authentication. Expected \"%s\", got \"%s\"",
		message, knsign, signed)
	if signed == knsign {
		return params, true
	} else {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "Invalid signed token",
			},
		)
		return params, false
	}
}

func (self *HTTPServer) AllPrices(c *gin.Context) {
	log.Printf("Getting all prices \n")
	data, err := self.app.GetAllPrices(getTimePoint(c, true))
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
	log.Printf("Getting price for %s - %s \n", base, quote)
	pair, err := common.NewTokenPair(base, quote)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": "Token pair is not supported"},
		)
	} else {
		data, err := self.app.GetOnePrice(pair.PairID(), getTimePoint(c, true))
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

func (self *HTTPServer) AuthData(c *gin.Context) {
	log.Printf("Getting current auth data snapshot \n")
	_, ok := self.Authenticated(c, []string{})
	if !ok {
		return
	}

	data, err := self.app.GetAuthData(getTimePoint(c, true))
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
	log.Printf("Getting all rates \n")
	data, err := self.app.GetAllRates(getTimePoint(c, true))
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
	postForm, ok := self.Authenticated(c, []string{"sources", "dests", "rates", "expires"})
	if !ok {
		return
	}

	sources := postForm.Get("sources")
	dests := postForm.Get("dests")
	rates := postForm.Get("rates")
	blocks := postForm.Get("expiries")
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
	id, err := self.core.SetRates(sourceTokens, destTokens, bigRates, expiryBlocks)
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
				"id":      id,
			},
		)
	}
}

func (self *HTTPServer) Trade(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{"base", "quote", "amount", "rate", "type"})
	if !ok {
		return
	}

	exchangeParam := c.Param("exchangeid")
	baseTokenParam := postForm.Get("base")
	quoteTokenParam := postForm.Get("quote")
	amountParam := postForm.Get("amount")
	rateParam := postForm.Get("rate")
	typeParam := postForm.Get("type")

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
	log.Printf("http server: Trade: rate: %f, raw rate: %s", rate, rateParam)
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
	id, done, remaining, finished, err := self.core.Trade(
		exchange, typeParam, base, quote, rate, amount, getTimePoint(c, false))
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
			"id":        id,
			"done":      done,
			"remaining": remaining,
			"finished":  finished,
		},
	)
}

func (self *HTTPServer) CancelOrder(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{"order_id"})
	if !ok {
		return
	}

	exchangeParam := c.Param("exchangeid")
	id := postForm.Get("order_id")

	exchange, err := common.GetExchange(exchangeParam)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	log.Printf("Cancel order id: %s from %s\n", id, exchange.ID())
	activityID, err := common.StringToActivityID(id)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": err.Error()},
		)
		return
	}
	err = self.core.CancelOrder(activityID, exchange)
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

func (self *HTTPServer) Withdraw(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{"token", "amount"})
	if !ok {
		return
	}

	exchangeParam := c.Param("exchangeid")
	tokenParam := postForm.Get("token")
	amountParam := postForm.Get("amount")

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
	log.Printf("Withdraw %s %s from %s\n", amount.Text(10), token.ID, exchange.ID())
	id, err := self.core.Withdraw(exchange, token, amount, getTimePoint(c, false))
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
			"id":      id,
		},
	)
}

func (self *HTTPServer) Deposit(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{"amount", "token"})
	if !ok {
		return
	}

	exchangeParam := c.Param("exchangeid")
	amountParam := postForm.Get("amount")
	tokenParam := postForm.Get("token")

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
	log.Printf("Depositing %s %s to %s\n", amount.Text(10), token.ID, exchange.ID())
	id, err := self.core.Deposit(exchange, token, amount, getTimePoint(c, false))
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
			"id":      id,
		},
	)
}

func (self *HTTPServer) GetActivities(c *gin.Context) {
	log.Printf("Getting all activity records \n")
	_, ok := self.Authenticated(c, []string{})
	if !ok {
		return
	}

	data, err := self.app.GetRecords()
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

func (self *HTTPServer) StopFetcher(c *gin.Context) {
	err := self.app.Stop()
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
			},
		)
	}
}

func (self *HTTPServer) ImmediatePendingActivities(c *gin.Context) {
	log.Printf("Getting all immediate pending activity records \n")
	_, ok := self.Authenticated(c, []string{})
	if !ok {
		return
	}

	data, err := self.app.GetPendingActivities()
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
	self.r.GET("/getrates", self.GetRate)

	self.r.GET("/authdata", self.AuthData)
	self.r.GET("/activities", self.GetActivities)
	self.r.GET("/immediate-pending-activities", self.ImmediatePendingActivities)

	self.r.POST("/cancelorder/:exchangeid", self.CancelOrder)
	self.r.POST("/deposit/:exchangeid", self.Deposit)
	self.r.POST("/withdraw/:exchangeid", self.Withdraw)
	self.r.POST("/trade/:exchangeid", self.Trade)
	self.r.POST("/setrates", self.SetRate)

	self.r.Run(self.host)
}

func NewHTTPServer(
	app reserve.ReserveData,
	core reserve.ReserveCore,
	host string,
	enableAuth bool,
	authEngine Authentication) *HTTPServer {
	raven.SetDSN("https://bf15053001464a5195a81bc41b644751:eff41ac715114b20b940010208271b13@sentry.io/228067")

	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.Use(cors.Default())

	return &HTTPServer{
		app, core, host, enableAuth, authEngine, r,
	}
}
