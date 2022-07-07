package servers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareHealthcheck(t *testing.T) {
	tests := map[string]struct {
		LiveStatus              bool
		ReadyStatus             bool
		expectedLivenessStatus  int
		expectedReadinessStatus int
	}{
		"everything down": {
			LiveStatus:              false,
			ReadyStatus:             false,
			expectedLivenessStatus:  http.StatusServiceUnavailable,
			expectedReadinessStatus: http.StatusServiceUnavailable,
		},
		"application down and service up": {
			LiveStatus:              false,
			ReadyStatus:             true,
			expectedLivenessStatus:  http.StatusServiceUnavailable,
			expectedReadinessStatus: http.StatusServiceUnavailable,
		},
		"application up and service down": {
			LiveStatus:              true,
			ReadyStatus:             false,
			expectedLivenessStatus:  http.StatusOK,
			expectedReadinessStatus: http.StatusServiceUnavailable,
		},
		"everything up": {
			LiveStatus:              true,
			ReadyStatus:             true,
			expectedLivenessStatus:  http.StatusOK,
			expectedReadinessStatus: http.StatusOK,
		},
	}

	handler := PrepareHealthcheck().Handler
	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	client := testServer.Client()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.LiveStatus {
				LiveStatus.MarkAsUp()
			} else {
				LiveStatus.MarkAsDown()
			}

			if tc.ReadyStatus {
				ReadyStatus.MarkAsUp()
			} else {
				ReadyStatus.MarkAsDown()
			}

			liveness, _ := client.Get(testServer.URL + "/live")
			readiness, _ := client.Get(testServer.URL + "/ready")

			require.Equal(t, tc.expectedLivenessStatus, liveness.StatusCode, "unexpected liveness status")
			require.Equal(t, tc.expectedReadinessStatus, readiness.StatusCode, "unexpected readiness status")
		})
	}
}
