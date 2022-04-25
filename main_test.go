package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_server(t *testing.T) {
	testcase := []struct {
		name   string
		header http.Header
		expect []byte
	}{
		{
			name:   "with-robot-forwarded-header",
			header: http.Header{"X-Forwarded-For": []string{"192.168.1.1"}, "User-Agent": []string{"curl/"}},
			expect: []byte("192.168.1.1\n"),
		},
		{
			name:   "with-robot-ua",
			header: http.Header{"User-Agent": []string{"curl/"}},
			expect: []byte("127.0.0.1\n"),
		},
	}

	for _, test := range testcase {
		require := require.New(t)
		mux := http.NewServeMux()
		handlers(mux)
		svr := httptest.NewServer(mux)
		c := svr.Client()
		req, err := http.NewRequest(http.MethodGet, svr.URL, nil)
		require.NoError(err)

		for k, v := range test.header {
			req.Header[k] = v
		}

		resp, err := c.Do(req)
		require.NoError(err)

		b, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(test.expect, b)
	}
}

var result = `Headers:
	Abc: [123]
	This-Is-A: [test]
`

func Test_handler(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("abc", "123")
	req.Header.Set("this-is-a", "test")
	require.NoError(t, err)
	h := handler()
	h.ServeHTTP(recorder, req)

	require.Equal(t, recorder.Code, http.StatusOK)
	require.Equal(t, result, recorder.Body.String())
}
