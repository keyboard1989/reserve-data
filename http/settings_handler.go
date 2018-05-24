package http

import (
	"fmt"
	"log"
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

func removeTokenFromList(tokens []common.Token, t common.Token) ([]common.Token, error) {
	if len(tokens) < 1 {
		return tokens, fmt.Errorf("Internal Token list is empty")
	}
	for i, token := range tokens {
		if token.ID == t.ID {
			tokens[len(tokens)-1], tokens[i] = tokens[i], tokens[len(tokens)-1]
			return tokens[:len(tokens)-1], nil
		}
	}
	return tokens, fmt.Errorf("The deactivating token is not in current internal token list")
}

func (self *HTTPServer) ReloadTokenIndices(newToken common.Token, active bool) error {
	tokens, err := settings.GetInternalTokens()
	if err != nil {
		return err
	}
	if active {
		tokens = append(tokens, newToken)
	} else {
		log.Printf("TOkens before removal: %v", tokens)
		tokens, err = removeTokenFromList(tokens, newToken)
		log.Printf("TOkens after removal: %v", tokens)
		if err != nil {
			return err
		}
	}
	if err = self.blockchain.LoadAndSetTokenIndices(common.GetTokenAddressesList(tokens)); err != nil {
		return err
	}
	return nil
}

func (self *HTTPServer) UpdateToken(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{"name", "address", "id", "decimals", "internal", "listed"}, []Permission{RebalancePermission, ConfigurePermission})
	if !ok {
		return
	}
	addrs := postForm.Get("address")
	name := postForm.Get("name")
	ID := postForm.Get("id")
	decimal := postForm.Get("decimals")
	internal := postForm.Get("internal")
	active := postForm.Get("listed")
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
	//We only concern with internal token, of which we must query its Indices
	if internalBool {
		if err = self.ReloadTokenIndices(token, activeBool); err != nil {
			returnError(c, err)
			return
		}
	}
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
