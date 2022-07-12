package servers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartServer(t *testing.T) {
	testsArithmetic := map[string]struct {
		operation        string
		number1, number2 float64
		expectedResult   float64
		expectedStatus   int
		arg1             string
		arg2             string
	}{

		"test sum": {
			operation:      "sum",
			number1:        3,
			number2:        4,
			expectedResult: 7,
			expectedStatus: http.StatusOK,
		},
		"test diff": {
			operation:      "diff",
			number1:        5,
			number2:        4,
			expectedResult: 1,
			expectedStatus: http.StatusOK,
		},
		"test mul": {
			operation:      "mul",
			number1:        7,
			number2:        4,
			expectedResult: 28,
			expectedStatus: http.StatusOK,
		},
		"test div": {
			operation:      "div",
			number1:        6,
			number2:        4,
			expectedResult: 1.5,
			expectedStatus: http.StatusOK,
		},
		"test fact": {
			operation:      "fact",
			number1:        3,
			expectedResult: 6,
			expectedStatus: http.StatusOK,
		},
		"test edge sum": {
			operation:      "sum",
			number1:        5.5,
			number2:        3.7,
			expectedResult: 9.2,
			expectedStatus: http.StatusOK,
		},
		"test edge diff": {
			operation:      "diff",
			number1:        6.6,
			number2:        4.1,
			expectedResult: 2.5,
			expectedStatus: http.StatusOK,
		},
		"test edge mul": {
			operation:      "mul",
			number1:        2.5,
			number2:        2.5,
			expectedResult: 6.25,
			expectedStatus: http.StatusOK,
		},
		"test edge fact": {
			operation:      "fact",
			number1:        2.5,
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test edge div": {
			operation:      "div",
			number1:        7.5,
			number2:        3,
			expectedResult: 2.5,
			expectedStatus: http.StatusOK,
		},
		"test div by 0": {
			operation:      "div",
			number1:        7.7,
			number2:        0,
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong args sum": {
			operation:      "sum",
			arg1:           "asd",
			arg2:           "xyz",
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong args diff": {
			operation:      "diff",
			arg1:           "asd",
			arg2:           "xyz",
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong args mul": {
			operation:      "mul",
			arg1:           "asd",
			arg2:           "xyz",
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong args div": {
			operation:      "div",
			arg1:           "asd",
			arg2:           "xyz",
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong args fact": {
			operation:      "fact",
			arg1:           "asd",
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong arg1 ": {
			operation:      "sum",
			arg1:           "asd",
			number2:        6,
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
		"test wrong arg2": {
			operation:      "sum",
			number1:        8,
			arg2:           "xyz",
			expectedResult: 0,
			expectedStatus: http.StatusBadRequest,
		},
	}
	handler := PrepareServer().Handler
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	client := testServer.Client()

	for name, tc := range testsArithmetic {

		t.Run(name, func(t *testing.T) {
			var status *http.Response
			if tc.operation == "fact" {
				var a string
				value := (int)(tc.number1)
				if (float64)(value) == tc.number1 {
					a = fmt.Sprintf("%d", value)
				} else {
					a = fmt.Sprintf("%f", tc.number1)
				}
				if tc.arg1 != "" {
					a = tc.arg1
				}
				status, _ = client.Get(testServer.URL + "/" + tc.operation + "/" + a)

			} else {
				var a = fmt.Sprintf("%f", tc.number1)
				var b = fmt.Sprintf("%f", tc.number2)
				if tc.arg1 != "" {
					a = tc.arg1
				}
				if tc.arg2 != "" {
					b = tc.arg1
				}
				status, _ = client.Get(testServer.URL + "/" + tc.operation + "/" + a + "/" + b)
			}

			bodyBytes, _ := ioutil.ReadAll(status.Body)
			bodyString := string(bodyBytes)
			result, _ := strconv.ParseFloat(bodyString, 64)
			require.Equal(t, tc.expectedStatus, status.StatusCode, "unexpected status")
			require.Equal(t, tc.expectedResult, result, "unexpected result")

		})

	}

}
