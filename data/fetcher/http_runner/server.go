package http_runner

import (
	"github.com/gin-gonic/gin"
	"os"
	"log"
	"github.com/gin-contrib/sentry"
	"github.com/getsentry/raven-go"
	"time"
	"net/http"
)

type HttpRunnerServer struct {
	runner HttpRunner
	host   string
	r      *gin.Engine
}

func (self *HttpRunnerServer) etick(c *gin.Context) {
	timestamp := time.Now()
	self.runner.eticker <- timestamp
	c.JSON(
		http.StatusOK,
		gin.H{
			"success":   true,
			"timestamp": timestamp,
		},
	)
}

func (self *HttpRunnerServer) btick(c *gin.Context) {
	timestamp := time.Now()
	self.runner.bticker <- timestamp
	c.JSON(
		http.StatusOK,
		gin.H{
			"success":   true,
			"timestamp": timestamp,
		},
	)
}

func (self *HttpRunnerServer) Run() {
	self.r.GET("/etick", self.etick)
	self.r.GET("/btick", self.btick)
	self.r.Run(self.host)
}

func NewHttpRunnerServer(runner HttpRunner, host string) *HttpRunnerServer {
	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	return &HttpRunnerServer{
		runner,
		host,
		r,
	}
}
