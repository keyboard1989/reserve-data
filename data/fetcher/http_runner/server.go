package http_runner

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sentry"
	"github.com/getsentry/raven-go"
	"time"
	"net/http"
	"errors"
)

type HttpRunnerServer struct {
	runner *HttpRunner
	host   string
	r      *gin.Engine
	http   *http.Server
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

func (self *HttpRunnerServer) init() {
	self.r.GET("/etick", self.etick)
	self.r.GET("/btick", self.btick)
}

func (self *HttpRunnerServer) Start() error {
	if self.http == nil {
		self.http = &http.Server{
			Addr:    self.host,
			Handler: self.r,
		}
		return self.http.ListenAndServe()
	} else {
		return errors.New("server start already")
	}
}

func (self *HttpRunnerServer) Stop() error {
	if self.http != nil {
		err := self.http.Shutdown(nil)
		self.http = nil
		return err
	} else {
		return errors.New("server stop already")
	}
}

func NewHttpRunnerServer(runner *HttpRunner, host string) *HttpRunnerServer {
	r := gin.Default()
	r.Use(sentry.Recovery(raven.DefaultClient, false))
	server := HttpRunnerServer{
		runner,
		host,
		r,
		nil,
	}
	server.init()
	return &server
}
