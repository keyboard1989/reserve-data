package http

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

func removeTokenFromList(tokens []common.Token, t common.Token) ([]common.Token, error) {
	if len(tokens) == 0 {
		return tokens, errors.New("Internal Token list is empty")
	}
	for i, token := range tokens {
		if token.ID == t.ID {
			tokens[len(tokens)-1], tokens[i] = tokens[i], tokens[len(tokens)-1]
			return tokens[:len(tokens)-1], nil
		}
	}
	return tokens, fmt.Errorf("The deactivating token %s is not in current internal token list", t.ID)
}

func (self *HTTPServer) reloadTokenIndices(newToken common.Token, active bool) error {
	tokens, err := self.setting.GetInternalTokens()
	if err != nil {
		return err
	}
	if active {
		tokens = append(tokens, newToken)
	} else {
		if tokens, err = removeTokenFromList(tokens, newToken); err != nil {
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
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	activeBool, err := strconv.ParseBool(active)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	internalBool, err := strconv.ParseBool(internal)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	token := common.NewToken(ID, name, addrs, decimalint64, activeBool, internalBool, minrr, maxpbi, maxti)
	//We only concern reloading indices if the token is Internal
	if internalBool {
		if err = self.reloadTokenIndices(token, activeBool); err != nil {
			httputil.ResponseFailure(c, httputil.WithError(err))
			return
		}
	}
	err = self.setting.UpdateToken(token)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

func (self *HTTPServer) TokenSettings(c *gin.Context) {
	_, ok := self.Authenticated(c, []string{}, []Permission{RebalancePermission, ConfigurePermission, ReadOnlyPermission, ConfirmConfPermission})
	if !ok {
		return
	}
	data, err := self.setting.GetAllTokens()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
	return
}
