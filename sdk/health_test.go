package sdk

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_AddHealthcheck(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	h := &Healthcheck{}
	require.NoError(h.AddHealthcheck("ok", HealthcheckerFunc(AlwaysOK)))
	require.NoError(h.AddHealthcheck("not-ok", HealthcheckerFunc(AlwaysFailing)))

	for name, check := range h.checks {
		result := runCheck(context.Background(), name, check)
		if name == "ok" {
			require.Equal(result.Result, 0)
			require.Equal(result.Name, "ok")
		} else {
			require.Equal(result.Result, 1)
			require.Equal(result.Name, "not-ok")
		}
	}
}

func Test_ParallelAddHealthcheck(t *testing.T) {
	h := &Healthcheck{}
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("n%d", i), func(t *testing.T) {
			t.Parallel()
			h.AddHealthcheck("ok", HealthcheckerFunc(AlwaysOK))
			h.AddHealthcheck("ok", HealthcheckerFunc(AlwaysFailing))
		})
	}
}

func Test_Serve(t *testing.T) {
	testcases := []struct {
		name         string
		responseCode int
		extraCheck   Healthchecker
	}{
		{
			name:         "ok",
			responseCode: http.StatusOK,
			extraCheck:   nil,
		},
		{
			name:         "timeout",
			responseCode: http.StatusGatewayTimeout,
			extraCheck: HealthcheckerFunc(func(ctx context.Context) error {
				time.Sleep(2 * time.Second)
				return nil
			}),
		},
	}

	healthcheckTimeout = 1 * time.Second
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			h := &Healthcheck{}

			h.AddHealthcheck("ok", HealthcheckerFunc(AlwaysOK))
			if tc.extraCheck != nil {
				h.AddHealthcheck("timeout", tc.extraCheck)
			}

			srv := httptest.NewServer(h)
			resp, err := http.Get(srv.URL)
			r.NoError(err)
			r.Equal(tc.responseCode, resp.StatusCode)

		})
	}
}
