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
	postForm, ok := self.Authenticated(c, []string{"address", "ID", "decimals", "internal", "active"}, []Permission{RebalancePermission, ConfigurePermission})
	if !ok {
		return
	}
	addrs := postForm.Get("address")
	ID := postForm.Get("ID")
	decimal := postForm.Get("decimals")
	internal := postForm.Get("internal")
	active := postForm.Get("active")

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

	token := common.NewToken(ID, addrs, decimalint64, activeBool, internalBool)

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
