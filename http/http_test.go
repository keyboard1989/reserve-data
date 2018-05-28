package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"net/url"

	"os"

	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func TestHTTPServerPWIEquationV2(t *testing.T) {
	const (
		endpoint = "/v2/set-pwis-equation"
		testData = `{
  "EOS": {
    "bid": {
      "a": 750,
      "b": 500,
      "c": 0,
      "MinMinSpread": 0,
      "PriceMultiplyFactor": 0
    },
    "ask": {
      "a": 800,
      "b": 600,
      "c": 0,
      "MinMinSpread": 0,
      "PriceMultiplyFactor": 0
    }
  },
  "ETH": {
    "bid": {
      "a": 750,
      "b": 500,
      "c": 0,
      "MinMinSpread": 0,
      "PriceMultiplyFactor": 0
    },
    "ask": {
      "a": 800,
      "b": 600,
      "c": 0,
      "MinMinSpread": 0,
      "PriceMultiplyFactor": 0
    }
  }
}
`
	)
	tmpDir, err := ioutil.TempDir("", "test_pwi_equation_v2")
	if err != nil {
		t.Fatal(err)
	}
	st, err := storage.NewBoltStorage(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}

	s := HTTPServer{
		app:         data.NewReserveData(st, nil, nil, nil, nil),
		core:        core.NewReserveCore(nil, st, common.Address{}),
		metric:      st,
		authEnabled: false,
		r:           gin.Default()}
	s.register()

	var tests = []struct {
		msg             string
		data            string
		expectedSuccess bool
	}{
		{
			msg:             "testing invalid post form",
			data:            `{"invalid_key": "invalid_value"}`,
			expectedSuccess: false,
		},
		{
			msg:             "testing valid post form",
			data:            testData,
			expectedSuccess: true,
		},
		// setting PWI equations when pending exists are not allowed
		{
			msg:             "testing setting when pending exists",
			data:            testData,
			expectedSuccess: false,
		},
	}

	for _, tc := range tests {
		t.Log(tc.msg)

		req, tErr := http.NewRequest(http.MethodPost, endpoint, nil)
		if tErr != nil {
			t.Fatal(tErr)
		}

		form := url.Values{}
		form.Add("data", tc.data)
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp := httptest.NewRecorder()
		s.r.ServeHTTP(resp, req)

		if tc.expectedSuccess {
			httputil.ExpectSuccess(t, resp)
		} else {
			httputil.ExpectFailure(t, resp)
		}
	}

	if err = os.RemoveAll(tmpDir); err != nil {
		t.Error(err)
	}
}
