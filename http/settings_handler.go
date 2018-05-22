package http

import (
	"net/http"
	"strconv"

	"github.com/KyberNetwork/reserve-data/settings"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/gin-gonic/gin"
)

func returnError(c *gin.Context, err error) {

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": false,
			"reason":  err.Error(),
		},
	)
	return
}

func (self *HTTPServer) UpdateToken(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{"name", "address", "id", "decimals", "internal", "active", "minimalRecordResolution", "maxPerBlockImbalance", "maxTotalImbalance"}, []Permission{RebalancePermission, ConfigurePermission})
	if !ok {
		return
	}
	addrs := postForm.Get("address")
	name := postForm.Get("name")
	ID := postForm.Get("id")
	decimal := postForm.Get("decimals")
	internal := postForm.Get("internal")
	active := postForm.Get("active")
	minrr := postForm.Get("minimalRecordResolution")
	maxpbi := postForm.Get("maxPerBlockImbalance")
	maxti := postForm.Get("maxTotalImbalance")
	decimalint64, err := strconv.ParseInt(decimal, 10, 64)
	if err != nil {
		returnError(c, err)
		return
	}
	activeBool, err := strconv.ParseBool(active)
	if err != nil {
		returnError(c, err)
		return
	}
	internalBool, err := strconv.ParseBool(internal)
	if err != nil {
		returnError(c, err)
		return
	}
	token := common.NewToken(ID, name, addrs, decimalint64, activeBool, internalBool, minrr, maxpbi, maxti)

	err = settings.UpdateToken(token)
	if err != nil {
		returnError(c, err)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
		},
	)
}

func (self *HTTPServer) TokenSettings(c *gin.Context) {
	_, ok := self.Authenticated(c, []string{}, []Permission{RebalancePermission, ConfigurePermission, ReadOnlyPermission, ConfirmConfPermission})
	if !ok {
		return
	}
	data, err := settings.GetAllTokens()
	if err != nil {
		returnError(c, err)
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success":   true,
			"timestamp": common.GetTimepoint(),
			"data":      data,
		},
	)
	return
}
