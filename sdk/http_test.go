package sdk

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewMuxWithHealthcheck(t *testing.T) {
	require := require.New(t)
	mux := NewMux()
	svc := httptest.NewServer(mux)

	c := svc.Client()
	u, err := url.Parse(svc.URL)
	require.NoError(err)
	u.Path = "/healthz"
	resp, err := c.Get(u.String())
	require.NoError(err)
	require.Equal(resp.StatusCode, 200)
}

/*
func Test_NewMuxWithOverriddeHealthcheck(t *testing.T) {
	require := require.New(t)
	mux := NewMux()
	svc := httptest.NewServer(mux)
	mux.Handle("/healthz", healthz(AlwaysFailing))

	c := svc.Client()
	url, err := url.Parse(svc.URL)
	require.NoError(err)
	url.Path = "/healthz"

	resp, err := c.Get(url.String())
	require.NoError(err)
	require.Equal(resp.StatusCode, 512)
}
*/
