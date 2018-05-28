package http

import (
	"encoding/json"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/metric"
	"github.com/gin-gonic/gin"
)

// SetPWIEquationV2 stores the given PWI equations to pending for later evaluation.
func (self *HTTPServer) SetPWIEquationV2(c *gin.Context) {
	postForm, ok := self.Authenticated(c, []string{}, []Permission{ConfigurePermission})
	if !ok {
		return
	}

	data := []byte(postForm.Get("data"))
	if len(data) > MAX_DATA_SIZE {
		httputil.ResponseFailure(c, httputil.WithError(errDataSizeExceed))
		return
	}

	var input metric.PWIEquationRequestV2
	if err := json.Unmarshal(data, &input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	if !input.IsValid() {
		httputil.ResponseFailure(c, httputil.WithReason("invalid input"))
		return
	}

	if err := self.metric.StorePendingPWIEquationV2(data); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

// GetPendingPWIEquationV2 returns the pending PWI equations.
func (self *HTTPServer) GetPendingPWIEquationV2(c *gin.Context) {
	_, ok := self.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, RebalancePermission, ConfigurePermission, ConfirmConfPermission})
	if !ok {
		return
	}
	data, err := self.metric.GetPendingPWIEquationV2()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}
