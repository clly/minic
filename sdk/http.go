package sdk

import (
	"errors"
	"net/http"
	"time"
)

var httpSvr *http.Server

const addr = "0.0.0.0:8080"
const readWriteTimeout = 30 * time.Second
const headerTimeout = 5 * time.Second

//const maxLength = 1024 ^ 1024*5

func ConfigureHTTP(handler http.Handler) *http.Server {
	httpSvr = &http.Server{
		Addr:              addr,
		Handler:           handler,
		TLSConfig:         nil,
		ReadTimeout:       readWriteTimeout,
		ReadHeaderTimeout: headerTimeout,
		WriteTimeout:      readWriteTimeout,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	return httpSvr
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz(AlwaysOK))
	return mux
}

func AlwaysOK() error { return nil }

func AlwaysFailing() error { return errors.New("oopsie") }

func healthz(health func() error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := health()
		if err != nil {
			w.WriteHeader(512)
		} else {
			w.WriteHeader(200)
		}
	})
}
