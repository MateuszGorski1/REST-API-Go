package servers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartServer(t *testing.T) {
	tests := map[string]struct {
		expectedSumStatus        int
		expectedDiffStatus       int
		expectedMulStatus        int
		expectedDivStatus        int
		expectedFactStatus       int
		expectedBadRequestStatus int
	}{

		"everything up": {
			expectedSumStatus:        http.StatusOK,
			expectedDiffStatus:       http.StatusOK,
			expectedMulStatus:        http.StatusOK,
			expectedDivStatus:        http.StatusOK,
			expectedFactStatus:       http.StatusOK,
			expectedBadRequestStatus: http.StatusBadRequest,
		},
	}

	handler := PrepareServer().Handler
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	client := testServer.Client()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			sumStatus, _ := client.Get(testServer.URL + "/sum/0/0")
			diffStatus, _ := client.Get(testServer.URL + "/diff/0/0")
			mulStatus, _ := client.Get(testServer.URL + "/mul/0/0")
			divStatus, _ := client.Get(testServer.URL + "/div/0/0")
			factStatus, _ := client.Get(testServer.URL + "/fact/1")
			badRequestStatus, _ := client.Get(testServer.URL + "/sum")

			require.Equal(t, tc.expectedSumStatus, sumStatus.StatusCode, "unexpected sum status")
			require.Equal(t, tc.expectedDiffStatus, diffStatus.StatusCode, "unexpected diff status")
			require.Equal(t, tc.expectedMulStatus, mulStatus.StatusCode, "unexpected mul status")
			require.Equal(t, tc.expectedDivStatus, divStatus.StatusCode, "unexpected div status")
			require.Equal(t, tc.expectedFactStatus, factStatus.StatusCode, "unexpected fact status")
			require.Equal(t, tc.expectedBadRequestStatus, badRequestStatus.StatusCode, "unexpected badrequest status")
		})
	}
}
