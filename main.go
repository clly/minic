package main

import (
	"fmt"
	"go.clly.me/minic/sdk"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("hello world")
	mux := serveMux()

	srv := sdk.ConfigureHTTP(mux)

	// TODO: Easier signal handling
	log.Fatal(srv.ListenAndServe())
}

func serveMux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", handler())
	return m
}

func handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if isRobot(request.UserAgent()) {
			remoteAddr := request.Header.Get("X-Remote-Addr")
			if remoteAddr == "" {
				remoteAddr = strings.Split(request.RemoteAddr, ":")[0]
			}
			fmt.Fprintln(writer, remoteAddr)
			return
		}
		fmt.Fprintln(writer, "Headers:")
		for k, v := range request.Header {
			fmt.Fprintf(writer, "\t%s: %s\n", k, v)
		}
	}
}

var robotRegexp = regexp.MustCompile(`^curl/.*$`)

func isRobot(ua string) bool {
	return robotRegexp.MatchString(ua)
}
