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
	postForm, ok := self.Authenticated(c, []string{"address", "ID", "decimals", "internal", "active", "minimalRecordResolution", "maxPerBlockImbalance", "maxTotalImbalance"}, []Permission{RebalancePermission, ConfigurePermission})
	if !ok {
		return
	}
	addrs := postForm.Get("address")
	ID := postForm.Get("ID")
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
	minrrUint64, err := strconv.ParseUint(minrr, 10, 64)
	if err != nil {
		returnError(c, err)
	}
	maxpbiUint64, err := strconv.ParseUint(maxpbi, 10, 64)
	if err != nil {
		returnError(c, err)
	}
	maxtiUint64, err := strconv.ParseUint(maxti, 10, 64)
	if err != nil {
		returnError(c, err)
	}
	token := common.NewToken(ID, addrs, decimalint64, activeBool, internalBool, minrrUint64, maxpbiUint64, maxtiUint64)

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
