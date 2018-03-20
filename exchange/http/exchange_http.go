package http

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app         Huobi
	host        string
	authEnabled bool
	auth        Authentication
	r           *gin.Engine
}

func IsIntime(nonce string) bool {
	serverTime := common.GetTimepoint()
	log.Printf("Server time: %d, None: %d", serverTime, nonce)
	nonceInt, err := strconv.ParseInt(nonce, 10, 64)
	if err != nil {
		log.Printf("IsIntime returns false, err: %v", err)
		return false
	}
	difference := nonceInt - int64(serverTime)
	if difference < -30000 || difference > 30000 {
		log.Printf("IsIntime returns false, nonce: %d, serverTime: %d, difference: %d", nonceInt, int64(serverTime), difference)
		return false
	}
	return true
}

func eligible(ups, allowedPerms []Permission) bool {
	for _, up := range ups {
		for _, ap := range allowedPerms {
			if up == ap {
				return true
			}
		}
	}
	return false
}

func (self *HTTPServer) Authenticated(c *gin.Context, requiredParams []string, perms []Permission) (url.Values, bool) {
	//parse the request
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
	//if no authentication, return true.
	if !self.authEnabled {
		return c.Request.Form, true
	}
	params := c.Request.Form
	log.Printf("Form params: %s\n", params)
	//if nonce is invalid
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
	//Check params
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
	//take sign message and check against user permission
	signed := c.GetHeader("signed")
	message := c.Request.Form.Encode()
	userPerms := self.auth.GetPermission(signed, message)
	if eligible(userPerms, perms) {
		return params, true
	} else {
		if len(userPerms) == 0 {
			c.JSON(
				http.StatusOK,
				gin.H{
					"success": false,
					"reason":  "Invalid signed token",
				},
			)
		} else {
			c.JSON(
				http.StatusOK,
				gin.H{
					"success": false,
					"reason":  "You don't have permission to proceed",
				},
			)
		}
		return params, false
	}
}

func (self *HTTPServer) PendingIntermediateTxs(c *gin.Context) {
	_, ok := self.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, RebalancePermission, ConfigurePermission, ConfirmConfPermission})
	if !ok {
		return
	}
	data, err := self.app.PendingIntermediateTxs()
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
	if self.app != nil {
		self.r.GET("/pending_intermediate_tx", self.PendingIntermediateTxs)
	}

	self.r.Run(self.host)
}

func NewHuobiHTTPServer(
	app Huobi,
	host string,
	enableAuth bool,
	authEngine Authentication,
	env string) *HTTPServer {

	r := gin.Default()
	sentryCli, err := raven.NewWithTags(
		"https://bf15053001464a5195a81bc41b644751:eff41ac715114b20b940010208271b13@sentry.io/228067",
		map[string]string{
			"env": env,
		},
	)
	if err != nil {
		panic(err)
	}
	r.Use(sentry.Recovery(
		sentryCli,
		false,
	))
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("signed")
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute
	r.Use(cors.New(corsConfig))

	return &HTTPServer{
		app, host, enableAuth, authEngine, r,
	}
}
