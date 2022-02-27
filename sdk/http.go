package sdk

import (
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
