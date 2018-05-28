package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/metric"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

type testCase struct {
	msg             string
	endpoint        string
	method          string
	data            string
	expectedSuccess bool
	assert          assertFn
}

func TestHTTPServerPWIEquationV2(t *testing.T) {
	const (
		storePendingPWIEquationV2Endpoint = "/v2/set-pwis-equation"
		getPendingPWIEquationV2Endpoint   = "/v2/pending-pwis-equation"
		testData                          = `{
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

	var tests = []testCase{
		{
			msg:             "testing invalid post form",
			endpoint:        storePendingPWIEquationV2Endpoint,
			method:          http.MethodPost,
			data:            `{"invalid_key": "invalid_value"}`,
			expectedSuccess: false,
		},
		{
			msg:             "getting non exists pending equation",
			endpoint:        getPendingPWIEquationV2Endpoint,
			method:          http.MethodGet,
			expectedSuccess: false,
		},
		{
			msg:             "testing valid post form",
			endpoint:        storePendingPWIEquationV2Endpoint,
			method:          http.MethodPost,
			data:            testData,
			expectedSuccess: true,
		},
		{
			msg:             "testing setting when pending exists",
			endpoint:        storePendingPWIEquationV2Endpoint,
			method:          http.MethodPost,
			data:            testData,
			expectedSuccess: false,
		},
		{
			msg:      "getting existing pending equation",
			endpoint: getPendingPWIEquationV2Endpoint,
			method:   http.MethodGet,
			assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				if resp.Code != http.StatusOK {
					t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
				}

				type responseBody struct {
					Success bool
					Data    metric.PWIEquationRequestV2
				}

				decoded := responseBody{}
				if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
					t.Fatal(err)
				}

				if decoded.Success != true {
					t.Errorf("wrong success status, expected: %t, got: %t", true, decoded.Success)
				}

				t.Logf("returned pending PWI equation request: %v", decoded.Data)

				if len(decoded.Data) != 2 {
					t.Fatalf("wrong number of tokens, expected: %d, got %d", 2, len(decoded.Data))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Log(tc.msg)

		req, tErr := http.NewRequest(tc.method, tc.endpoint, nil)
		if tErr != nil {
			t.Fatal(tErr)
		}

		if tc.data != "" {
			form := url.Values{}
			form.Add("data", tc.data)
			req.PostForm = form
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}

		resp := httptest.NewRecorder()
		s.r.ServeHTTP(resp, req)

		if tc.assert == nil {
			if tc.expectedSuccess {
				tc.assert = httputil.ExpectSuccess
			} else {
				tc.assert = httputil.ExpectFailure
			}
		}
		tc.assert(t, resp)
	}

	if err = os.RemoveAll(tmpDir); err != nil {
		t.Error(err)
	}
}
