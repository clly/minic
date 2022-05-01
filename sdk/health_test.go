package sdk

import (
	"context"
	"testing"

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
