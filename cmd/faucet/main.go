package main

import (
	"fmt"
	"math/big"
	"net/http"

	ethereum "github.com/ethereum/go-ethereum/common"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

type FaucetServer struct {
	r   *gin.Engine
	app *FaucetApp
}

func (self *FaucetServer) Claim(c *gin.Context) {
	addr := ethereum.HexToAddress(c.PostForm("address"))
	if addr.Big().Cmp(big.NewInt(0)) == 0 {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "error": "Invalid address"},
		)
	} else {
		sent, found := self.app.Get(addr)
		if found {
			c.JSON(
				http.StatusOK,
				gin.H{
					"success": false,
					"error":   fmt.Sprintf("Your address is already registered. We have sent ETH to your address with tx: %s", sent.Hex()),
					"tx":      sent.Hex(),
				},
			)
		} else {
			no, added := self.app.AddAddress(addr)
			if !added {
				latestIndex, yourIndex, _ := self.app.Search(addr)
				c.JSON(
					http.StatusOK,
					gin.H{
						"success": false,
						"error":   fmt.Sprintf("Your address is already registered. If you haven't receive the ETH yet, please wait, there are %d addresses before yours.", yourIndex-latestIndex)},
				)
			} else {
				c.JSON(
					http.StatusOK,
					gin.H{"success": true, "msg": fmt.Sprintf("Your address is added to faucet queue. There are %d addresses before yours. Please wait.", no)},
				)
			}
		}
	}
}

func (self *FaucetServer) Run() {
	self.r.POST("/claim-eth", self.Claim)
	go self.app.Run()
	self.r.Run(":8891")
}

func main() {
	raven.SetDSN("https://bf15053001464a5195a81bc41b644751:eff41ac715114b20b940010208271b13@sentry.io/228067")
	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.Use(cors.Default())

	server := &FaucetServer{
		r, NewFaucetApp(),
	}

	server.Run()
}
